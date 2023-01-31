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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/edgefarm/provider-nats/apis/consumer/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeConsumers implements ConsumerInterface
type FakeConsumers struct {
	Fake *FakeConsumerV1alpha1
}

var consumersResource = schema.GroupVersionResource{Group: "consumer", Version: "v1alpha1", Resource: "consumers"}

var consumersKind = schema.GroupVersionKind{Group: "consumer", Version: "v1alpha1", Kind: "Consumer"}

// Get takes name of the consumer, and returns the corresponding consumer object, and an error if there is any.
func (c *FakeConsumers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Consumer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(consumersResource, name), &v1alpha1.Consumer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Consumer), err
}

// List takes label and field selectors, and returns the list of Consumers that match those selectors.
func (c *FakeConsumers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ConsumerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(consumersResource, consumersKind, opts), &v1alpha1.ConsumerList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ConsumerList{ListMeta: obj.(*v1alpha1.ConsumerList).ListMeta}
	for _, item := range obj.(*v1alpha1.ConsumerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested consumers.
func (c *FakeConsumers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(consumersResource, opts))
}

// Create takes the representation of a consumer and creates it.  Returns the server's representation of the consumer, and an error, if there is any.
func (c *FakeConsumers) Create(ctx context.Context, consumer *v1alpha1.Consumer, opts v1.CreateOptions) (result *v1alpha1.Consumer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(consumersResource, consumer), &v1alpha1.Consumer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Consumer), err
}

// Update takes the representation of a consumer and updates it. Returns the server's representation of the consumer, and an error, if there is any.
func (c *FakeConsumers) Update(ctx context.Context, consumer *v1alpha1.Consumer, opts v1.UpdateOptions) (result *v1alpha1.Consumer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(consumersResource, consumer), &v1alpha1.Consumer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Consumer), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeConsumers) UpdateStatus(ctx context.Context, consumer *v1alpha1.Consumer, opts v1.UpdateOptions) (*v1alpha1.Consumer, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(consumersResource, "status", consumer), &v1alpha1.Consumer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Consumer), err
}

// Delete takes name of the consumer and deletes it. Returns an error if one occurs.
func (c *FakeConsumers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(consumersResource, name, opts), &v1alpha1.Consumer{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeConsumers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(consumersResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ConsumerList{})
	return err
}

// Patch applies the patch and returns the patched consumer.
func (c *FakeConsumers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Consumer, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(consumersResource, name, pt, data, subresources...), &v1alpha1.Consumer{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Consumer), err
}