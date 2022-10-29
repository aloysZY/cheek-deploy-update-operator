//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022 aloys.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheekDeployUpdate) DeepCopyInto(out *CheekDeployUpdate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheekDeployUpdate.
func (in *CheekDeployUpdate) DeepCopy() *CheekDeployUpdate {
	if in == nil {
		return nil
	}
	out := new(CheekDeployUpdate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CheekDeployUpdate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheekDeployUpdateList) DeepCopyInto(out *CheekDeployUpdateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]CheekDeployUpdate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheekDeployUpdateList.
func (in *CheekDeployUpdateList) DeepCopy() *CheekDeployUpdateList {
	if in == nil {
		return nil
	}
	out := new(CheekDeployUpdateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CheekDeployUpdateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheekDeployUpdateSpec) DeepCopyInto(out *CheekDeployUpdateSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheekDeployUpdateSpec.
func (in *CheekDeployUpdateSpec) DeepCopy() *CheekDeployUpdateSpec {
	if in == nil {
		return nil
	}
	out := new(CheekDeployUpdateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheekDeployUpdateStatus) DeepCopyInto(out *CheekDeployUpdateStatus) {
	*out = *in
	in.CDUStatus.DeepCopyInto(&out.CDUStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheekDeployUpdateStatus.
func (in *CheekDeployUpdateStatus) DeepCopy() *CheekDeployUpdateStatus {
	if in == nil {
		return nil
	}
	out := new(CheekDeployUpdateStatus)
	in.DeepCopyInto(out)
	return out
}
