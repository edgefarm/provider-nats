package stream_test

import (
	"encoding/json"
	"fmt"

	v1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	natsgo "github.com/nats-io/nats.go"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	streamsv1 "github.com/edgefarm/provider-nats/apis/stream/v1alpha1"
	utils "github.com/edgefarm/provider-nats/cluster/local/e2e/pkg/utils"
	"github.com/edgefarm/provider-nats/internal/clients/nats"
)

//nolint:unparam
func userCreds(f *utils.Framework, secret string) (string, string, string, error) {
	acc, err := f.GetSecret("crossplane-system", secret)
	if err != nil {
		return "", "", "", err
	}
	var x map[string]interface{}
	j, ok := acc["credentials"]
	if !ok {
		return "", "", "", fmt.Errorf("credentials not found in secret")
	}
	err = json.Unmarshal(j, &x)
	if err != nil {
		return "", "", "", err
	}
	fmt.Fprintf(GinkgoWriter, "x: %+v", x)
	return x["jwt"].(string), x["seed_key"].(string), x["address"].(string), nil
}

func pubKey(f *utils.Framework, secret string) (string, string, error) {
	acc, err := f.GetSecret("crossplane-system", secret)
	if err != nil {
		return "", "", err
	}
	var x map[string]interface{}
	j, ok := acc["credentials"]
	if !ok {
		return "", "", fmt.Errorf("credentials not found in secret")
	}
	err = json.Unmarshal(j, &x)
	if err != nil {
		return "", "", err
	}
	fmt.Fprintf(GinkgoWriter, "x: %+v", x)
	apub, upub, err := nats.GetPublicKeys(x["jwt"].(string))
	if err != nil {
		return "", "", err
	}
	return apub, upub, nil
}

var _ = Describe("Stream E2E tests", func() {
	Context("Handle streams for main nats (no domain given)", func() {
		var (
			f *utils.Framework
		)
		BeforeEach(func() {
			f = utils.DefaultFramework
		})
		Context("Stream main-acc1", func() {
			It("can create the stream", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/main-acc1.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateOrLeaveStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())
			})
			It("can update the stream", func() {
				streamOld, err := f.GetStream("main-acc1")
				Expect(err).To(BeNil())
				streamNew := streamOld.DeepCopy()
				streamNew.Spec.ForProvider.Config.MaxBytes = 402400
				streamNew.Spec.ForProvider.Config.Discard = "New"
				err = f.UpdateStream(streamNew)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(streamNew.ObjectMeta.Name)
				Expect(err).To(BeNil())

				Expect(streamOld.Spec.ForProvider.Config.MaxBytes).Should(BeIdenticalTo(int64(102400)))
				Expect(streamNew.Spec.ForProvider.Config.MaxBytes).Should(BeIdenticalTo(int64(402400)))
				Expect(streamOld.Spec.ForProvider.Config.Discard).Should(BeIdenticalTo("Old"))
				Expect(streamNew.Spec.ForProvider.Config.Discard).Should(BeIdenticalTo("New"))
			})
		})
		Context("Stream main-acc2", func() {
			It("should succeed on creating main-acc2 for the first time", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/main-acc2.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())
			})
			It("should fail creating main-acc2 on the second time", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/main-acc2.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(kerrors.IsAlreadyExists(err)).To(BeTrue())
			})
		})
		Context("Validate stream connection info", func() {
			It("can validate for main-acc1", func() {
				domain, err := f.GetStreamDomain("main-acc1")
				Expect(err).To(BeNil())
				Expect(domain).Should(Equal(""))

				accountPubKey, userPubKey, err := pubKey(f, "acc1-creds")
				Expect(err).To(BeNil())
				acc1UserPubKey, err := f.GetStreamUserPublicKey("main-acc1")
				Expect(err).To(BeNil())
				Expect(acc1UserPubKey).Should(Equal(userPubKey))
				acc1PubKey, err := f.GetStreamAccountPublicKey("main-acc1")
				Expect(err).To(BeNil())
				Expect(acc1PubKey).Should(Equal(accountPubKey))
			})
			It("can validate for main-acc2", func() {
				domain, err := f.GetStreamDomain("main-acc2")
				Expect(err).To(BeNil())
				Expect(domain).Should(Equal(""))

				accountPubKey, userPubKey, err := pubKey(f, "acc2-creds")
				Expect(err).To(BeNil())
				acc1UserPubKey, err := f.GetStreamUserPublicKey("main-acc2")
				Expect(err).To(BeNil())
				Expect(acc1UserPubKey).Should(Equal(userPubKey))
				acc1PubKey, err := f.GetStreamAccountPublicKey("main-acc2")
				Expect(err).To(BeNil())
				Expect(acc1PubKey).Should(Equal(accountPubKey))
			})
		})
		Context("Delete streams", func() {
			It("can delete main-acc1", func() {
				err := f.DeleteStream("main-acc1")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("main-acc1")
				Expect(err).To(BeNil())
			})
			It("can delete main-acc2", func() {
				err := f.DeleteStream("main-acc2")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("main-acc2")
				Expect(err).To(BeNil())
			})
		})
	})
	Context("Handle streams for leaf nats servers (domain given)", func() {
		var (
			f *utils.Framework
		)
		BeforeEach(func() {
			f = utils.DefaultFramework
		})
		Context("Create leaf nats streams", func() {
			It("can create the stream for domain foo", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/foo-acc1.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())
			})
			It("can create the stream for domain bar", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/bar-acc2.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())
			})
			It("can't create the stream for domain baz", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/std/baz-acc1.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamConditionsAvailable(stream.ObjectMeta.Name, 2)
				Expect(err).To(BeNil())
				current, err := f.GetStream(stream.ObjectMeta.Name)
				conditionTime := current.Status.Conditions[0].LastTransitionTime
				Expect(err).To(BeNil())
				Expect(v1.Condition{
					Type:               v1.TypeReady,
					Status:             corev1.ConditionFalse,
					LastTransitionTime: conditionTime,
					Reason:             v1.ReasonUnavailable,
					Message:            natsgo.ErrNoResponders.Error(),
				}).Should(BeElementOf(current.Status.Conditions))
			})
		})
		Context("Delete leaf nats streams", func() {
			It("can delete the stream for domain foo", func() {
				err := f.DeleteStream("foo-acc1")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("foo-acc1")
				Expect(err).To(BeNil())
			})
			It("can delete the stream for domain bar", func() {
				err := f.DeleteStream("bar-acc2")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("foo-acc1")
				Expect(err).To(BeNil())
			})
			It("can delete the stream for domain baz", func() {
				err := f.DeleteStream("baz-acc1")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("foo-acc1")
				Expect(err).To(BeNil())
			})
		})
	})
	Context("Handle sourced streams", func() {
		var (
			f *utils.Framework
		)
		BeforeEach(func() {
			f = utils.DefaultFramework
		})
		Context("Create sourced stream", func() {
			It("can create the leaf node streams", func() {
				streamFoo := &streamsv1.Stream{}
				streamFoo.Spec.ForProvider.Config.SetDefaults()
				streamFoo, err := utils.UnmarshalAnyYaml("manifests/stream/aggregate/source_foo.yaml", streamFoo)
				Expect(err).To(BeNil())
				err = f.CreateStream(streamFoo)
				Expect(err).To(BeNil())

				streamBoth := &streamsv1.Stream{}
				streamBoth.Spec.ForProvider.Config.SetDefaults()
				streamBoth, err = utils.UnmarshalAnyYaml("manifests/stream/aggregate/source_both.yaml", streamBoth)
				Expect(err).To(BeNil())
				err = f.CreateStream(streamBoth)
				Expect(err).To(BeNil())
			})
			It("can create the aggregating stream", func() {
				stream := &streamsv1.Stream{}
				stream.Spec.ForProvider.Config.SetDefaults()
				stream, err := utils.UnmarshalAnyYaml("manifests/stream/aggregate/aggregate.yaml", stream)
				Expect(err).To(BeNil())
				err = f.CreateStream(stream)
				Expect(err).To(BeNil())
				err = f.WaitForStreamSyncAndReady(stream.ObjectMeta.Name)
				Expect(err).To(BeNil())
				parameters, err := f.GetStreamParameters("aggregate")
				Expect(err).To(BeNil())
				Expect(parameters.Config.Sources).Should(HaveLen(3))
			})
		})
		// Sending 10 messages should result in 20 messages in the aggregate stream because it
		// is sends to the domains `foo` and `both` have a stream that is subscribed
		// to the same subject. The aggregate stream sources from both domains and therefore should
		// receive 20 messages.
		Context("Stream gets synced", func() {
			It("can publish messages", func() {
				jwt, seed, address, err := userCreds(f, "acc1-creds")
				Expect(err).To(BeNil())
				p := &utils.Pub{
					JWT:     jwt,
					Seed:    seed,
					Address: address,
					Subject: "source.foo",
					Count:   "10",
					Data:    "foo",
				}
				err = f.NatsPublish(p)
				Expect(err).To(BeNil())
			})
			It("can sync `aggregate`", func() {
				jwt, seed, address, err := userCreds(f, "acc1-creds")
				Expect(err).To(BeNil())
				err = f.NatsCheckMessages(&utils.StreamInfo{
					JWT:      jwt,
					Seed:     seed,
					Address:  address,
					Count:    "20",
					Stream:   "aggregate",
					Timeout:  "20",
					Interval: "1",
					Domain:   "",
				})
				Expect(err).To(BeNil())
			})
			It("can sync `foo-source`", func() {
				jwt, seed, address, err := userCreds(f, "acc1-creds")
				Expect(err).To(BeNil())
				err = f.NatsCheckMessages(&utils.StreamInfo{
					JWT:      jwt,
					Seed:     seed,
					Address:  address,
					Count:    "10",
					Stream:   "source",
					Timeout:  "20",
					Interval: "1",
					Domain:   "foo",
				})
				Expect(err).To(BeNil())
			})
			It("can sync `both-source`", func() {
				jwt, seed, address, err := userCreds(f, "acc1-creds")
				Expect(err).To(BeNil())
				err = f.NatsCheckMessages(&utils.StreamInfo{
					JWT:      jwt,
					Seed:     seed,
					Address:  address,
					Count:    "10",
					Stream:   "source",
					Timeout:  "20",
					Interval: "1",
					Domain:   "both",
				})
				Expect(err).To(BeNil())
			})
		})
		Context("Stream status gets updated", func() {
			It("can update status for `aggregate`", func() {
				err := f.WaitForStreamMessags("aggregate", 20)
				Expect(err).To(BeNil())
			})
			It("can update status for `foo-source`", func() {
				err := f.WaitForStreamMessags("foo-source", 10)
				Expect(err).To(BeNil())
			})
			It("can update status for `both-source`", func() {
				err := f.WaitForStreamMessags("both-source", 10)
				Expect(err).To(BeNil())
			})
		})
		Context("Delete sourced stream", func() {
			It("can delete `aggregate`", func() {
				err := f.DeleteStream("aggregate")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("aggregate")
				Expect(err).To(BeNil())
			})
			It("can delete `source` on domain `foo`", func() {
				err := f.DeleteStream("foo-source")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("foo-source")
				Expect(err).To(BeNil())
			})
			It("can delete `source` on domain `both`", func() {
				err := f.DeleteStream("both-source")
				Expect(err).To(BeNil())
				err = f.WaitForStreamDeleted("both-source")
				Expect(err).To(BeNil())
			})
		})
	})
})
