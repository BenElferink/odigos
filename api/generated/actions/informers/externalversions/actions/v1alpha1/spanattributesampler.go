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
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"
	time "time"

	apiactionsv1alpha1 "github.com/odigos-io/odigos/api/actions/v1alpha1"
	versioned "github.com/odigos-io/odigos/api/generated/actions/clientset/versioned"
	internalinterfaces "github.com/odigos-io/odigos/api/generated/actions/informers/externalversions/internalinterfaces"
	actionsv1alpha1 "github.com/odigos-io/odigos/api/generated/actions/listers/actions/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SpanAttributeSamplerInformer provides access to a shared informer and lister for
// SpanAttributeSamplers.
type SpanAttributeSamplerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() actionsv1alpha1.SpanAttributeSamplerLister
}

type spanAttributeSamplerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSpanAttributeSamplerInformer constructs a new informer for SpanAttributeSampler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSpanAttributeSamplerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSpanAttributeSamplerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSpanAttributeSamplerInformer constructs a new informer for SpanAttributeSampler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSpanAttributeSamplerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ActionsV1alpha1().SpanAttributeSamplers(namespace).List(context.Background(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ActionsV1alpha1().SpanAttributeSamplers(namespace).Watch(context.Background(), options)
			},
			ListWithContextFunc: func(ctx context.Context, options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ActionsV1alpha1().SpanAttributeSamplers(namespace).List(ctx, options)
			},
			WatchFuncWithContext: func(ctx context.Context, options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ActionsV1alpha1().SpanAttributeSamplers(namespace).Watch(ctx, options)
			},
		},
		&apiactionsv1alpha1.SpanAttributeSampler{},
		resyncPeriod,
		indexers,
	)
}

func (f *spanAttributeSamplerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSpanAttributeSamplerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *spanAttributeSamplerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiactionsv1alpha1.SpanAttributeSampler{}, f.defaultInformer)
}

func (f *spanAttributeSamplerInformer) Lister() actionsv1alpha1.SpanAttributeSamplerLister {
	return actionsv1alpha1.NewSpanAttributeSamplerLister(f.Informer().GetIndexer())
}
