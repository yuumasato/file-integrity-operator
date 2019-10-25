// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileIntegrity) DeepCopyInto(out *FileIntegrity) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileIntegrity.
func (in *FileIntegrity) DeepCopy() *FileIntegrity {
	if in == nil {
		return nil
	}
	out := new(FileIntegrity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FileIntegrity) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileIntegrityConfig) DeepCopyInto(out *FileIntegrityConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileIntegrityConfig.
func (in *FileIntegrityConfig) DeepCopy() *FileIntegrityConfig {
	if in == nil {
		return nil
	}
	out := new(FileIntegrityConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileIntegrityList) DeepCopyInto(out *FileIntegrityList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FileIntegrity, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileIntegrityList.
func (in *FileIntegrityList) DeepCopy() *FileIntegrityList {
	if in == nil {
		return nil
	}
	out := new(FileIntegrityList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FileIntegrityList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileIntegritySpec) DeepCopyInto(out *FileIntegritySpec) {
	*out = *in
	out.Config = in.Config
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileIntegritySpec.
func (in *FileIntegritySpec) DeepCopy() *FileIntegritySpec {
	if in == nil {
		return nil
	}
	out := new(FileIntegritySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FileIntegrityStatus) DeepCopyInto(out *FileIntegrityStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FileIntegrityStatus.
func (in *FileIntegrityStatus) DeepCopy() *FileIntegrityStatus {
	if in == nil {
		return nil
	}
	out := new(FileIntegrityStatus)
	in.DeepCopyInto(out)
	return out
}