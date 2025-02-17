// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: mathsoolset.proto

package sgc7pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RunScript - run script
type RunScript struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Script        string                 `protobuf:"bytes,1,opt,name=script,proto3" json:"script,omitempty"`
	MapFiles      string                 `protobuf:"bytes,2,opt,name=mapFiles,proto3" json:"mapFiles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RunScript) Reset() {
	*x = RunScript{}
	mi := &file_mathsoolset_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RunScript) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunScript) ProtoMessage() {}

func (x *RunScript) ProtoReflect() protoreflect.Message {
	mi := &file_mathsoolset_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunScript.ProtoReflect.Descriptor instead.
func (*RunScript) Descriptor() ([]byte, []int) {
	return file_mathsoolset_proto_rawDescGZIP(), []int{0}
}

func (x *RunScript) GetScript() string {
	if x != nil {
		return x.Script
	}
	return ""
}

func (x *RunScript) GetMapFiles() string {
	if x != nil {
		return x.MapFiles
	}
	return ""
}

// ReplyRunScript - reply run script
type ReplyRunScript struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ScriptErrs    []string               `protobuf:"bytes,1,rep,name=scriptErrs,proto3" json:"scriptErrs,omitempty"`
	MapFiles      string                 `protobuf:"bytes,2,opt,name=mapFiles,proto3" json:"mapFiles,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReplyRunScript) Reset() {
	*x = ReplyRunScript{}
	mi := &file_mathsoolset_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReplyRunScript) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyRunScript) ProtoMessage() {}

func (x *ReplyRunScript) ProtoReflect() protoreflect.Message {
	mi := &file_mathsoolset_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyRunScript.ProtoReflect.Descriptor instead.
func (*ReplyRunScript) Descriptor() ([]byte, []int) {
	return file_mathsoolset_proto_rawDescGZIP(), []int{1}
}

func (x *ReplyRunScript) GetScriptErrs() []string {
	if x != nil {
		return x.ScriptErrs
	}
	return nil
}

func (x *ReplyRunScript) GetMapFiles() string {
	if x != nil {
		return x.MapFiles
	}
	return ""
}

var File_mathsoolset_proto protoreflect.FileDescriptor

var file_mathsoolset_proto_rawDesc = string([]byte{
	0x0a, 0x11, 0x6d, 0x61, 0x74, 0x68, 0x73, 0x6f, 0x6f, 0x6c, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x73, 0x67, 0x63, 0x37, 0x70, 0x62, 0x22, 0x3f, 0x0a, 0x09, 0x52,
	0x75, 0x6e, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x61, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6d, 0x61, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x4c, 0x0a, 0x0e,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x52, 0x75, 0x6e, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x45, 0x72, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0a, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x45, 0x72, 0x72, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x6d, 0x61, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6d, 0x61, 0x70, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x32, 0x47, 0x0a, 0x0b, 0x4d, 0x61,
	0x74, 0x68, 0x54, 0x6f, 0x6f, 0x6c, 0x73, 0x65, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x72, 0x75, 0x6e,
	0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x12, 0x11, 0x2e, 0x73, 0x67, 0x63, 0x37, 0x70, 0x62, 0x2e,
	0x52, 0x75, 0x6e, 0x53, 0x63, 0x72, 0x69, 0x70, 0x74, 0x1a, 0x16, 0x2e, 0x73, 0x67, 0x63, 0x37,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x52, 0x75, 0x6e, 0x53, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x22, 0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x7a, 0x68, 0x73, 0x30, 0x30, 0x37, 0x2f, 0x73, 0x6c, 0x6f, 0x74, 0x73, 0x67, 0x61,
	0x6d, 0x65, 0x63, 0x6f, 0x72, 0x65, 0x37, 0x2f, 0x73, 0x67, 0x63, 0x37, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_mathsoolset_proto_rawDescOnce sync.Once
	file_mathsoolset_proto_rawDescData []byte
)

func file_mathsoolset_proto_rawDescGZIP() []byte {
	file_mathsoolset_proto_rawDescOnce.Do(func() {
		file_mathsoolset_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_mathsoolset_proto_rawDesc), len(file_mathsoolset_proto_rawDesc)))
	})
	return file_mathsoolset_proto_rawDescData
}

var file_mathsoolset_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_mathsoolset_proto_goTypes = []any{
	(*RunScript)(nil),      // 0: sgc7pb.RunScript
	(*ReplyRunScript)(nil), // 1: sgc7pb.ReplyRunScript
}
var file_mathsoolset_proto_depIdxs = []int32{
	0, // 0: sgc7pb.MathToolset.runScript:input_type -> sgc7pb.RunScript
	1, // 1: sgc7pb.MathToolset.runScript:output_type -> sgc7pb.ReplyRunScript
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mathsoolset_proto_init() }
func file_mathsoolset_proto_init() {
	if File_mathsoolset_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_mathsoolset_proto_rawDesc), len(file_mathsoolset_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mathsoolset_proto_goTypes,
		DependencyIndexes: file_mathsoolset_proto_depIdxs,
		MessageInfos:      file_mathsoolset_proto_msgTypes,
	}.Build()
	File_mathsoolset_proto = out.File
	file_mathsoolset_proto_goTypes = nil
	file_mathsoolset_proto_depIdxs = nil
}
