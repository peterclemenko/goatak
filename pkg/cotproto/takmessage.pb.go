// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: takmessage.proto

package cotproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Top level message sent for TAK Messaging Protocol Version 1.
type TakMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional - if omitted, continue using last reported control
	// information
	TakControl *TakControl `protobuf:"bytes,1,opt,name=takControl,proto3" json:"takControl,omitempty"`
	// Optional - if omitted, no event data in this message
	CotEvent       *CotEvent `protobuf:"bytes,2,opt,name=cotEvent,proto3" json:"cotEvent,omitempty"`
	SubmissionTime uint64    `protobuf:"varint,3,opt,name=submissionTime,proto3" json:"submissionTime,omitempty"`
	CreationTime   uint64    `protobuf:"varint,4,opt,name=creationTime,proto3" json:"creationTime,omitempty"`
}

func (x *TakMessage) Reset() {
	*x = TakMessage{}

	if protoimpl.UnsafeEnabled {
		mi := &file_takmessage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TakMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TakMessage) ProtoMessage() {}

func (x *TakMessage) ProtoReflect() protoreflect.Message {
	mi := &file_takmessage_proto_msgTypes[0]

	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}

		return ms
	}

	return mi.MessageOf(x)
}

// Deprecated: Use TakMessage.ProtoReflect.Descriptor instead.
func (*TakMessage) Descriptor() ([]byte, []int) {
	return file_takmessage_proto_rawDescGZIP(), []int{0}
}

func (x *TakMessage) GetTakControl() *TakControl {
	if x != nil {
		return x.TakControl
	}

	return nil
}

func (x *TakMessage) GetCotEvent() *CotEvent {
	if x != nil {
		return x.CotEvent
	}

	return nil
}

func (x *TakMessage) GetSubmissionTime() uint64 {
	if x != nil {
		return x.SubmissionTime
	}

	return 0
}

func (x *TakMessage) GetCreationTime() uint64 {
	if x != nil {
		return x.CreationTime
	}

	return 0
}

var File_takmessage_proto protoreflect.FileDescriptor

var file_takmessage_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x61, 0x6b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0e, 0x63, 0x6f, 0x74, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x10, 0x74, 0x61, 0x6b, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x01, 0x0a, 0x0a, 0x54, 0x61, 0x6b, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x74, 0x61, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x54, 0x61, 0x6b, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x52, 0x0a, 0x74, 0x61, 0x6b, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x12, 0x25, 0x0a, 0x08, 0x63, 0x6f, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x09, 0x2e, 0x43, 0x6f, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x63,
	0x6f, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x0e, 0x73, 0x75, 0x62, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x0e, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x22, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x42, 0x26, 0x48, 0x03, 0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x64, 0x75, 0x64, 0x6b, 0x6f, 0x76, 0x2f, 0x67, 0x6f, 0x61, 0x74,
	0x61, 0x6b, 0x2f, 0x63, 0x6f, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_takmessage_proto_rawDescOnce sync.Once
	file_takmessage_proto_rawDescData = file_takmessage_proto_rawDesc
)

func file_takmessage_proto_rawDescGZIP() []byte {
	file_takmessage_proto_rawDescOnce.Do(func() {
		file_takmessage_proto_rawDescData = protoimpl.X.CompressGZIP(file_takmessage_proto_rawDescData)
	})

	return file_takmessage_proto_rawDescData
}

var file_takmessage_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_takmessage_proto_goTypes = []interface{}{
	(*TakMessage)(nil), // 0: TakMessage
	(*TakControl)(nil), // 1: TakControl
	(*CotEvent)(nil),   // 2: CotEvent
}
var file_takmessage_proto_depIdxs = []int32{
	1, // 0: TakMessage.takControl:type_name -> TakControl
	2, // 1: TakMessage.cotEvent:type_name -> CotEvent
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_takmessage_proto_init() }
func file_takmessage_proto_init() {
	if File_takmessage_proto != nil {
		return
	}

	file_cotevent_proto_init()
	file_takcontrol_proto_init()

	if !protoimpl.UnsafeEnabled {
		file_takmessage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TakMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_takmessage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_takmessage_proto_goTypes,
		DependencyIndexes: file_takmessage_proto_depIdxs,
		MessageInfos:      file_takmessage_proto_msgTypes,
	}.Build()
	File_takmessage_proto = out.File
	file_takmessage_proto_rawDesc = nil
	file_takmessage_proto_goTypes = nil
	file_takmessage_proto_depIdxs = nil
}
