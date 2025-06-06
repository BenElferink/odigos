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

	versioned "github.com/odigos-io/odigos/api/generated/odigos/clientset/versioned"
	internalinterfaces "github.com/odigos-io/odigos/api/generated/odigos/informers/externalversions/internalinterfaces"
	odigosv1alpha1 "github.com/odigos-io/odigos/api/generated/odigos/listers/odigos/v1alpha1"
	apiodigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CollectorsGroupInformer provides access to a shared informer and lister for
// CollectorsGroups.
type CollectorsGroupInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() odigosv1alpha1.CollectorsGroupLister
}

type collectorsGroupInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCollectorsGroupInformer constructs a new informer for CollectorsGroup type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCollectorsGroupInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCollectorsGroupInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCollectorsGroupInformer constructs a new informer for CollectorsGroup type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCollectorsGroupInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OdigosV1alpha1().CollectorsGroups(namespace).List(context.Background(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OdigosV1alpha1().CollectorsGroups(namespace).Watch(context.Background(), options)
			},
			ListWithContextFunc: func(ctx context.Context, options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OdigosV1alpha1().CollectorsGroups(namespace).List(ctx, options)
			},
			WatchFuncWithContext: func(ctx context.Context, options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.OdigosV1alpha1().CollectorsGroups(namespace).Watch(ctx, options)
			},
		},
		&apiodigosv1alpha1.CollectorsGroup{},
		resyncPeriod,
		indexers,
	)
}

func (f *collectorsGroupInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCollectorsGroupInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *collectorsGroupInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apiodigosv1alpha1.CollectorsGroup{}, f.defaultInformer)
}

func (f *collectorsGroupInformer) Lister() odigosv1alpha1.CollectorsGroupLister {
	return odigosv1alpha1.NewCollectorsGroupLister(f.Informer().GetIndexer())
}
