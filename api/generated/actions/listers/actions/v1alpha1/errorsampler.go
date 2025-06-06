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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	actionsv1alpha1 "github.com/odigos-io/odigos/api/actions/v1alpha1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// ErrorSamplerLister helps list ErrorSamplers.
// All objects returned here must be treated as read-only.
type ErrorSamplerLister interface {
	// List lists all ErrorSamplers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*actionsv1alpha1.ErrorSampler, err error)
	// ErrorSamplers returns an object that can list and get ErrorSamplers.
	ErrorSamplers(namespace string) ErrorSamplerNamespaceLister
	ErrorSamplerListerExpansion
}

// errorSamplerLister implements the ErrorSamplerLister interface.
type errorSamplerLister struct {
	listers.ResourceIndexer[*actionsv1alpha1.ErrorSampler]
}

// NewErrorSamplerLister returns a new ErrorSamplerLister.
func NewErrorSamplerLister(indexer cache.Indexer) ErrorSamplerLister {
	return &errorSamplerLister{listers.New[*actionsv1alpha1.ErrorSampler](indexer, actionsv1alpha1.Resource("errorsampler"))}
}

// ErrorSamplers returns an object that can list and get ErrorSamplers.
func (s *errorSamplerLister) ErrorSamplers(namespace string) ErrorSamplerNamespaceLister {
	return errorSamplerNamespaceLister{listers.NewNamespaced[*actionsv1alpha1.ErrorSampler](s.ResourceIndexer, namespace)}
}

// ErrorSamplerNamespaceLister helps list and get ErrorSamplers.
// All objects returned here must be treated as read-only.
type ErrorSamplerNamespaceLister interface {
	// List lists all ErrorSamplers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*actionsv1alpha1.ErrorSampler, err error)
	// Get retrieves the ErrorSampler from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*actionsv1alpha1.ErrorSampler, error)
	ErrorSamplerNamespaceListerExpansion
}

// errorSamplerNamespaceLister implements the ErrorSamplerNamespaceLister
// interface.
type errorSamplerNamespaceLister struct {
	listers.ResourceIndexer[*actionsv1alpha1.ErrorSampler]
}
