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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// OtelAttributeWithValueApplyConfiguration represents an declarative configuration of the OtelAttributeWithValue type for use
// with apply.
type OtelAttributeWithValueApplyConfiguration struct {
	AttributeName        *string `json:"attributeName,omitempty"`
	AttributeStringValue *string `json:"attributeStringValue,omitempty"`
}

// OtelAttributeWithValueApplyConfiguration constructs an declarative configuration of the OtelAttributeWithValue type for use with
// apply.
func OtelAttributeWithValue() *OtelAttributeWithValueApplyConfiguration {
	return &OtelAttributeWithValueApplyConfiguration{}
}

// WithAttributeName sets the AttributeName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AttributeName field is set to the value of the last call.
func (b *OtelAttributeWithValueApplyConfiguration) WithAttributeName(value string) *OtelAttributeWithValueApplyConfiguration {
	b.AttributeName = &value
	return b
}

// WithAttributeStringValue sets the AttributeStringValue field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AttributeStringValue field is set to the value of the last call.
func (b *OtelAttributeWithValueApplyConfiguration) WithAttributeStringValue(value string) *OtelAttributeWithValueApplyConfiguration {
	b.AttributeStringValue = &value
	return b
}