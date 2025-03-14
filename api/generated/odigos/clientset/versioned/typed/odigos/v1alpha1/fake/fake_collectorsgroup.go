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
	odigosv1alpha1 "github.com/odigos-io/odigos/api/generated/odigos/applyconfiguration/odigos/v1alpha1"
	typedodigosv1alpha1 "github.com/odigos-io/odigos/api/generated/odigos/clientset/versioned/typed/odigos/v1alpha1"
	v1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
	gentype "k8s.io/client-go/gentype"
)

// fakeCollectorsGroups implements CollectorsGroupInterface
type fakeCollectorsGroups struct {
	*gentype.FakeClientWithListAndApply[*v1alpha1.CollectorsGroup, *v1alpha1.CollectorsGroupList, *odigosv1alpha1.CollectorsGroupApplyConfiguration]
	Fake *FakeOdigosV1alpha1
}

func newFakeCollectorsGroups(fake *FakeOdigosV1alpha1, namespace string) typedodigosv1alpha1.CollectorsGroupInterface {
	return &fakeCollectorsGroups{
		gentype.NewFakeClientWithListAndApply[*v1alpha1.CollectorsGroup, *v1alpha1.CollectorsGroupList, *odigosv1alpha1.CollectorsGroupApplyConfiguration](
			fake.Fake,
			namespace,
			v1alpha1.SchemeGroupVersion.WithResource("collectorsgroups"),
			v1alpha1.SchemeGroupVersion.WithKind("CollectorsGroup"),
			func() *v1alpha1.CollectorsGroup { return &v1alpha1.CollectorsGroup{} },
			func() *v1alpha1.CollectorsGroupList { return &v1alpha1.CollectorsGroupList{} },
			func(dst, src *v1alpha1.CollectorsGroupList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.CollectorsGroupList) []*v1alpha1.CollectorsGroup {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.CollectorsGroupList, items []*v1alpha1.CollectorsGroup) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}
