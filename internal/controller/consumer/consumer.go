/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	natsgo "github.com/nats-io/nats.go"

	"github.com/edgefarm/provider-nats/apis/consumer/v1alpha1"
	"github.com/edgefarm/provider-nats/apis/consumer/v1alpha1/consumer"
	apisv1alpha1 "github.com/edgefarm/provider-nats/apis/v1alpha1"
	nats "github.com/edgefarm/provider-nats/internal/clients/nats"
	"github.com/edgefarm/provider-nats/internal/controller/features"
)

const (
	errNotConsumer  = "managed resource is not a Consumer custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

// Setup adds a controller that reconciles Consumer managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.ConsumerGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	connector := &connector{
		kube:         mgr.GetClient(),
		usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
		newServiceFn: nats.NewClient,
		logger:       o.Logger,
	}
	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ConsumerGroupVersionKind),
		managed.WithExternalConnecter(connector),
		managed.WithExternalConnectDisconnecter(connector),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&v1alpha1.Consumer{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*nats.Client, error)
	client       *nats.Client
	logger       logging.Logger
}

// Disconnect implements managed.ExternalConnectDisconnecter
func (c *connector) Disconnect(ctx context.Context) error {
	return c.client.Disconnect()
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Consumer)
	if !ok {
		return nil, errors.New(errNotConsumer)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	client, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	c.client = client
	e := &external{
		client:          client,
		clientCloseChan: make(chan struct{}),
		log:             c.logger,
	}

	return e, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	client          *nats.Client
	clientCloseChan chan struct{}
	log             logging.Logger
}

func (c *external) Disconnect(ctx context.Context) error {
	return c.client.Disconnect()
}

const (
	annotationExternalName = "crossplane.io/external-name"
)

func getExternalName(r *v1alpha1.Consumer) (string, error) {
	annotations := r.GetAnnotations()
	if annotations != nil {
		if val, ok := annotations[annotationExternalName]; ok {
			return val, nil
		}
	}
	return "", fmt.Errorf("External name annotation not found for stream %s", r.GetName())
}

func (c *external) setStatus(domain string, stream string, r *v1alpha1.Consumer, data *natsgo.ConsumerInfo) {
	r.Status.AtProvider.State.Domain = domain
	r.Status.AtProvider.State.Stream = stream
	r.Status.AtProvider.State.PushBound = func() string {
		if data.Config.DeliverSubject == "" {
			return "no"
		}
		return "yes"
	}()
	r.Status.AtProvider.State.Name = data.Config.Name
	r.Status.AtProvider.State.Durable = data.Config.Durable

	r.Status.AtProvider.State.Created = data.Created.String()
	r.Status.AtProvider.State.Delivered = consumer.SequenceInfo{
		Consumer: data.Delivered.Consumer,
		Stream:   data.Delivered.Stream,
	}
	if data.Delivered.Last != nil {
		r.Status.AtProvider.State.Delivered.Last = data.Delivered.Last.String()
	}

	r.Status.AtProvider.State.AckFloor = consumer.SequenceInfo{
		Consumer: data.AckFloor.Consumer,
		Stream:   data.AckFloor.Stream,
	}
	if data.AckFloor.Last != nil {
		r.Status.AtProvider.State.AckFloor.Last = data.AckFloor.Last.String()
	}

	r.Status.AtProvider.State.NumAckPending = data.NumAckPending
	r.Status.AtProvider.State.NumRedelivered = data.NumRedelivered
	r.Status.AtProvider.State.NumWaiting = data.NumWaiting
	r.Status.AtProvider.State.NumPending = data.NumPending
	if data.Cluster != nil {
		r.Status.AtProvider.State.Cluster = &consumer.ClusterInfo{
			Name:   data.Cluster.Name,
			Leader: data.Cluster.Leader,
		}
		if data.Cluster.Replicas != nil {
			for _, replica := range data.Cluster.Replicas {
				r.Status.AtProvider.State.Cluster.Replicas = append(r.Status.AtProvider.State.Cluster.Replicas, &consumer.PeerInfo{
					Name:    replica.Name,
					Current: replica.Current,
					Offline: replica.Offline,
					Active:  replica.Active.String(),
					Lag:     replica.Lag,
				})
			}
		}
	}
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	r, ok := mg.(*v1alpha1.Consumer)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotConsumer)
	}
	externalName, err := getExternalName(r)
	if err != nil {
		return managed.ExternalObservation{}, err
	}

	domain := r.Spec.ForProvider.Domain
	stream := r.Spec.ForProvider.Stream

	data, err := nats.ConsumerInfo(c.client, domain, externalName, stream)
	if err != nil {
		r.SetConditions(xpv1.Unavailable().WithMessage(err.Error()))
		return managed.ExternalObservation{}, err
	}

	if data == nil {
		r.SetConditions(xpv1.Unavailable())
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	customConfig := r.Spec.ForProvider.Config
	converted, err := consumer.ConfigV1Alpha1ToNats(externalName, &customConfig)
	if err != nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, err
	}

	converted.Name = externalName

	oriJson, err := json.Marshal(data.Config)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	convertedJson, err := json.Marshal(converted)
	if err != nil {
		return managed.ExternalObservation{}, err
	}
	// c.log.Debug("Diff", "upstream", oriJson, "converted", convertedJson)
	// fmt.Println("upst", string(oriJson))
	// fmt.Println("conv", string(convertedJson))

	if !bytes.Equal(oriJson, convertedJson) {
		return managed.ExternalObservation{
			ResourceExists:    true,
			ResourceUpToDate:  false,
			ConnectionDetails: managed.ConnectionDetails{},
		}, nil
	}

	c.setStatus(domain, stream, r, data)

	r.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		// Return false when the external resource does not exist. This lets
		// the managed resource reconciler know that it needs to call Create to
		// (re)create the resource, or that it has successfully been deleted.
		ResourceExists: true,

		// Return false when the external resource exists, but it not up to date
		// with the desired managed resource state. This lets the managed
		// resource reconciler know that it needs to call Update.
		ResourceUpToDate: true,

		// Return any details that may be required to connect to the external
		// resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	r, ok := mg.(*v1alpha1.Consumer)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotConsumer)
	}
	c.log.Info("Creating", "consumer", r)

	customConfig := r.Spec.ForProvider.Config
	domain := r.Spec.ForProvider.Domain
	stream := r.Spec.ForProvider.Stream
	externalName, err := getExternalName(r)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	config, err := consumer.ConfigV1Alpha1ToNats(externalName, &customConfig)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	config.Name = externalName
	err = c.client.CreateConsumer(domain, stream, config)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	return managed.ExternalCreation{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	r, ok := mg.(*v1alpha1.Consumer)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotConsumer)
	}

	c.log.Info("Updating", "consumer", r)

	customConfig := r.Spec.ForProvider.Config
	domain := r.Spec.ForProvider.Domain
	stream := r.Spec.ForProvider.Stream
	externalName, err := getExternalName(r)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	config, err := consumer.ConfigV1Alpha1ToNats(externalName, &customConfig)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}
	config.Name = externalName
	err = c.client.UpdateConsumer(domain, stream, config)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{
		// Optionally return any details that may be required to connect to the
		// external resource. These will be stored as the connection secret.
		ConnectionDetails: managed.ConnectionDetails{},
	}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	r, ok := mg.(*v1alpha1.Consumer)
	if !ok {
		return errors.New(errNotConsumer)
	}

	c.log.Info("Deleting", "consumer", r)

	domain := r.Spec.ForProvider.Domain
	stream := r.Spec.ForProvider.Stream
	externalName, err := getExternalName(r)
	if err != nil {
		return err
	}

	return c.client.DeleteConsumer(domain, stream, externalName)
}
