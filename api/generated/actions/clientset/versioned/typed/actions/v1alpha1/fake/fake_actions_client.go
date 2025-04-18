/*
Copyright 2022.

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
	v1alpha1 "github.com/odigos-io/odigos/api/generated/actions/clientset/versioned/typed/actions/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeActionsV1alpha1 struct {
	*testing.Fake
}

func (c *FakeActionsV1alpha1) AddClusterInfos(namespace string) v1alpha1.AddClusterInfoInterface {
	return newFakeAddClusterInfos(c, namespace)
}

func (c *FakeActionsV1alpha1) DeleteAttributes(namespace string) v1alpha1.DeleteAttributeInterface {
	return newFakeDeleteAttributes(c, namespace)
}

func (c *FakeActionsV1alpha1) ErrorSamplers(namespace string) v1alpha1.ErrorSamplerInterface {
	return newFakeErrorSamplers(c, namespace)
}

func (c *FakeActionsV1alpha1) K8sAttributesResolvers(namespace string) v1alpha1.K8sAttributesResolverInterface {
	return newFakeK8sAttributesResolvers(c, namespace)
}

func (c *FakeActionsV1alpha1) LatencySamplers(namespace string) v1alpha1.LatencySamplerInterface {
	return newFakeLatencySamplers(c, namespace)
}

func (c *FakeActionsV1alpha1) PiiMaskings(namespace string) v1alpha1.PiiMaskingInterface {
	return newFakePiiMaskings(c, namespace)
}

func (c *FakeActionsV1alpha1) ProbabilisticSamplers(namespace string) v1alpha1.ProbabilisticSamplerInterface {
	return newFakeProbabilisticSamplers(c, namespace)
}

func (c *FakeActionsV1alpha1) RenameAttributes(namespace string) v1alpha1.RenameAttributeInterface {
	return newFakeRenameAttributes(c, namespace)
}

func (c *FakeActionsV1alpha1) ServiceNameSamplers(namespace string) v1alpha1.ServiceNameSamplerInterface {
	return newFakeServiceNameSamplers(c, namespace)
}

func (c *FakeActionsV1alpha1) SpanAttributeSamplers(namespace string) v1alpha1.SpanAttributeSamplerInterface {
	return newFakeSpanAttributeSamplers(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeActionsV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
