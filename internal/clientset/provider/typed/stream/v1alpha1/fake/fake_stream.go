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

	v1alpha1 "github.com/edgefarm/provider-nats/apis/stream/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeStreams implements StreamInterface
type FakeStreams struct {
	Fake *FakeStreamV1alpha1
}

var streamsResource = schema.GroupVersionResource{Group: "stream", Version: "v1alpha1", Resource: "streams"}

var streamsKind = schema.GroupVersionKind{Group: "stream", Version: "v1alpha1", Kind: "Stream"}

// Get takes name of the stream, and returns the corresponding stream object, and an error if there is any.
func (c *FakeStreams) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Stream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(streamsResource, name), &v1alpha1.Stream{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stream), err
}

// List takes label and field selectors, and returns the list of Streams that match those selectors.
func (c *FakeStreams) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.StreamList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(streamsResource, streamsKind, opts), &v1alpha1.StreamList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.StreamList{ListMeta: obj.(*v1alpha1.StreamList).ListMeta}
	for _, item := range obj.(*v1alpha1.StreamList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested streams.
func (c *FakeStreams) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(streamsResource, opts))
}

// Create takes the representation of a stream and creates it.  Returns the server's representation of the stream, and an error, if there is any.
func (c *FakeStreams) Create(ctx context.Context, stream *v1alpha1.Stream, opts v1.CreateOptions) (result *v1alpha1.Stream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(streamsResource, stream), &v1alpha1.Stream{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stream), err
}

// Update takes the representation of a stream and updates it. Returns the server's representation of the stream, and an error, if there is any.
func (c *FakeStreams) Update(ctx context.Context, stream *v1alpha1.Stream, opts v1.UpdateOptions) (result *v1alpha1.Stream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(streamsResource, stream), &v1alpha1.Stream{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stream), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeStreams) UpdateStatus(ctx context.Context, stream *v1alpha1.Stream, opts v1.UpdateOptions) (*v1alpha1.Stream, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(streamsResource, "status", stream), &v1alpha1.Stream{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stream), err
}

// Delete takes name of the stream and deletes it. Returns an error if one occurs.
func (c *FakeStreams) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(streamsResource, name, opts), &v1alpha1.Stream{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStreams) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(streamsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.StreamList{})
	return err
}

// Patch applies the patch and returns the patched stream.
func (c *FakeStreams) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Stream, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(streamsResource, name, pt, data, subresources...), &v1alpha1.Stream{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Stream), err
}