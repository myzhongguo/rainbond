// RAINBOND, Application Management Platform
// Copyright (C) 2014-2020 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/goodrain/rainbond-operator/pkg/apis/rainbond/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeRbdComponents implements RbdComponentInterface
type FakeRbdComponents struct {
	Fake *FakeRainbondV1alpha1
	ns   string
}

var rbdcomponentsResource = schema.GroupVersionResource{Group: "rainbond.io", Version: "v1alpha1", Resource: "rbdcomponents"}

var rbdcomponentsKind = schema.GroupVersionKind{Group: "rainbond.io", Version: "v1alpha1", Kind: "RbdComponent"}

// Get takes name of the rbdComponent, and returns the corresponding rbdComponent object, and an error if there is any.
func (c *FakeRbdComponents) Get(name string, options v1.GetOptions) (result *v1alpha1.RbdComponent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(rbdcomponentsResource, c.ns, name), &v1alpha1.RbdComponent{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RbdComponent), err
}

// List takes label and field selectors, and returns the list of RbdComponents that match those selectors.
func (c *FakeRbdComponents) List(opts v1.ListOptions) (result *v1alpha1.RbdComponentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(rbdcomponentsResource, rbdcomponentsKind, c.ns, opts), &v1alpha1.RbdComponentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RbdComponentList{ListMeta: obj.(*v1alpha1.RbdComponentList).ListMeta}
	for _, item := range obj.(*v1alpha1.RbdComponentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested rbdComponents.
func (c *FakeRbdComponents) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(rbdcomponentsResource, c.ns, opts))

}

// Create takes the representation of a rbdComponent and creates it.  Returns the server's representation of the rbdComponent, and an error, if there is any.
func (c *FakeRbdComponents) Create(rbdComponent *v1alpha1.RbdComponent) (result *v1alpha1.RbdComponent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(rbdcomponentsResource, c.ns, rbdComponent), &v1alpha1.RbdComponent{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RbdComponent), err
}

// Update takes the representation of a rbdComponent and updates it. Returns the server's representation of the rbdComponent, and an error, if there is any.
func (c *FakeRbdComponents) Update(rbdComponent *v1alpha1.RbdComponent) (result *v1alpha1.RbdComponent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(rbdcomponentsResource, c.ns, rbdComponent), &v1alpha1.RbdComponent{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RbdComponent), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRbdComponents) UpdateStatus(rbdComponent *v1alpha1.RbdComponent) (*v1alpha1.RbdComponent, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(rbdcomponentsResource, "status", c.ns, rbdComponent), &v1alpha1.RbdComponent{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RbdComponent), err
}

// Delete takes name of the rbdComponent and deletes it. Returns an error if one occurs.
func (c *FakeRbdComponents) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(rbdcomponentsResource, c.ns, name), &v1alpha1.RbdComponent{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRbdComponents) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(rbdcomponentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RbdComponentList{})
	return err
}

// Patch applies the patch and returns the patched rbdComponent.
func (c *FakeRbdComponents) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RbdComponent, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(rbdcomponentsResource, c.ns, name, pt, data, subresources...), &v1alpha1.RbdComponent{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RbdComponent), err
}
