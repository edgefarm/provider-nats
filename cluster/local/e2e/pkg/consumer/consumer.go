package consumer_test

import (
	"fmt"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/nats-io/nats.go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	consumerv1 "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1"
	consumerinternalv1 "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1/consumer"
	errors "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1/consumer/errors"

	streamsv1 "github.com/edgefarm/provider-nats/apis/stream/v1alpha1"
	utils "github.com/edgefarm/provider-nats/cluster/local/e2e/pkg/utils"
)

var _ = Describe("Consumer E2E tests", func() {
	Context("Handle invalid consumer configurations", func() {
		var (
			f *utils.Framework
		)
		BeforeEach(func() {
			f = utils.DefaultFramework
		})
		Context("Conditions for consumer not met", func() {
			It("is missing the target stream", func() {
				consumer := &consumerv1.Consumer{}
				consumer.Spec.ForProvider.Config.SetDefaults(consumerinternalv1.ConsumerTypeNone)
				consumer, err := utils.UnmarshalAnyYaml("manifests/consumer/missing_stream.yaml", consumer)
				Expect(err).To(BeNil())
				err = f.CreateConsumer(consumer)
				Expect(err).To(BeNil())
				err = f.WaitForConsumerConditionsAvailable(consumer.ObjectMeta.Name, 2)
				Expect(err).To(BeNil())
				current, err := f.GetConsumer(consumer.ObjectMeta.Name)
				conditionTime := current.Status.Conditions[0].LastTransitionTime
				Expect(err).To(BeNil())
				Expect(v1.Condition{
					Type:               v1.TypeReady,
					Status:             corev1.ConditionFalse,
					LastTransitionTime: conditionTime,
					Reason:             v1.ReasonUnavailable,
					Message:            nats.ErrStreamNotFound.Error(),
				}).Should(BeElementOf(current.Status.Conditions))
			})
		})
		Context("Usage errors regarding consumer spec", func() {
			It("has both pull and push set", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/consumer/stream.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())

				consumer := &consumerv1.Consumer{}
				consumer.Spec.ForProvider.Config.SetDefaults(consumerinternalv1.ConsumerTypeNone)
				consumer, err = utils.UnmarshalAnyYaml("manifests/consumer/invalid_config.yaml", consumer)
				Expect(err).To(BeNil())
				err = f.CreateConsumer(consumer)
				Expect(err).To(BeNil())
				condition := &v1.Condition{
					Type:    v1.TypeSynced,
					Status:  corev1.ConditionFalse,
					Reason:  v1.ReasonReconcileError,
					Message: fmt.Sprintf("create failed: %s", errors.PushAndPullConsumerError.Error()),
				}
				err = f.WaitForConsumerCondition(consumer.ObjectMeta.Name, condition)
				Expect(err).To(BeNil())
			})
		})
	})
	Context("Handle invalid consumer configurations", func() {
		Context("Push and Pull consumers", func() {
			It("is can create the stream", func() {
				f := utils.DefaultFramework
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/consumer/stream.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateOrLeaveStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.Name)
				Expect(err).To(BeNil())
			})
			It("minimal consumer", func() {
				f := utils.DefaultFramework
				consumer := &consumerv1.Consumer{}
				consumer.Spec.ForProvider.Config.SetDefaults(consumerinternalv1.ConsumerTypeNone)
				consumer, err := utils.UnmarshalAnyYaml("manifests/consumer/minimal.yaml", consumer)
				Expect(err).To(BeNil())
				err = f.CreateConsumer(consumer)
				Expect(err).To(BeNil())
				err = f.WaitForConsumerSyncAndReady(consumer.Name)
				Expect(err).To(BeNil())
			})
			It("pull consumer", func() {
				f := utils.DefaultFramework
				consumer := &consumerv1.Consumer{}
				consumer.Spec.ForProvider.Config.SetDefaults(consumerinternalv1.ConsumerTypeNone)
				consumer, err := utils.UnmarshalAnyYaml("manifests/consumer/pull.yaml", consumer)
				Expect(err).To(BeNil())
				err = f.CreateConsumer(consumer)
				Expect(err).To(BeNil())
				err = f.WaitForConsumerSyncAndReady(consumer.Name)
				Expect(err).To(BeNil())
			})
			It("push consumer", func() {
				f := utils.DefaultFramework
				consumer := &consumerv1.Consumer{}
				consumer.Spec.ForProvider.Config.SetDefaults(consumerinternalv1.ConsumerTypeNone)
				consumer, err := utils.UnmarshalAnyYaml("manifests/consumer/push.yaml", consumer)
				Expect(err).To(BeNil())
				err = f.CreateConsumer(consumer)
				Expect(err).To(BeNil())
				err = f.WaitForConsumerSyncAndReady(consumer.Name)
				Expect(err).To(BeNil())
			})
		})
	})
})
