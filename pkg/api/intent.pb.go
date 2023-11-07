// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: proto/intent.proto

package api

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

type IntentType int32

const (
	IntentType_INTENT_TYPE_UNSPECIFIED     IntentType = 0
	IntentType_INTENT_TYPE_HIGH_BANDWIDTH  IntentType = 1
	IntentType_INTENT_TYPE_LOW_BANDWIDTH   IntentType = 2
	IntentType_INTENT_TYPE_LOW_LATENCY     IntentType = 3
	IntentType_INTENT_TYPE_LOW_PACKET_LOSS IntentType = 4
	IntentType_INTENT_TYPE_LOW_JITTER      IntentType = 5
)

// Enum value maps for IntentType.
var (
	IntentType_name = map[int32]string{
		0: "INTENT_TYPE_UNSPECIFIED",
		1: "INTENT_TYPE_HIGH_BANDWIDTH",
		2: "INTENT_TYPE_LOW_BANDWIDTH",
		3: "INTENT_TYPE_LOW_LATENCY",
		4: "INTENT_TYPE_LOW_PACKET_LOSS",
		5: "INTENT_TYPE_LOW_JITTER",
	}
	IntentType_value = map[string]int32{
		"INTENT_TYPE_UNSPECIFIED":     0,
		"INTENT_TYPE_HIGH_BANDWIDTH":  1,
		"INTENT_TYPE_LOW_BANDWIDTH":   2,
		"INTENT_TYPE_LOW_LATENCY":     3,
		"INTENT_TYPE_LOW_PACKET_LOSS": 4,
		"INTENT_TYPE_LOW_JITTER":      5,
	}
)

func (x IntentType) Enum() *IntentType {
	p := new(IntentType)
	*p = x
	return p
}

func (x IntentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IntentType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_intent_proto_enumTypes[0].Descriptor()
}

func (IntentType) Type() protoreflect.EnumType {
	return &file_proto_intent_proto_enumTypes[0]
}

func (x IntentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IntentType.Descriptor instead.
func (IntentType) EnumDescriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{0}
}

type ValueType int32

const (
	ValueType_VALUE_TYPE_UNSPECIFIED  ValueType = 0
	ValueType_VALUE_TYPE_MIN_VALUE    ValueType = 1
	ValueType_VALUE_TYPE_MAX_VALUE    ValueType = 2
	ValueType_VALUE_TYPE_SFC          ValueType = 3
	ValueType_VALUE_TYPE_FLEX_ALGO_NR ValueType = 4
)

// Enum value maps for ValueType.
var (
	ValueType_name = map[int32]string{
		0: "VALUE_TYPE_UNSPECIFIED",
		1: "VALUE_TYPE_MIN_VALUE",
		2: "VALUE_TYPE_MAX_VALUE",
		3: "VALUE_TYPE_SFC",
		4: "VALUE_TYPE_FLEX_ALGO_NR",
	}
	ValueType_value = map[string]int32{
		"VALUE_TYPE_UNSPECIFIED":  0,
		"VALUE_TYPE_MIN_VALUE":    1,
		"VALUE_TYPE_MAX_VALUE":    2,
		"VALUE_TYPE_SFC":          3,
		"VALUE_TYPE_FLEX_ALGO_NR": 4,
	}
)

func (x ValueType) Enum() *ValueType {
	p := new(ValueType)
	*p = x
	return p
}

func (x ValueType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ValueType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_intent_proto_enumTypes[1].Descriptor()
}

func (ValueType) Type() protoreflect.EnumType {
	return &file_proto_intent_proto_enumTypes[1]
}

func (x ValueType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ValueType.Descriptor instead.
func (ValueType) EnumDescriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{1}
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        ValueType `protobuf:"varint,1,opt,name=type,proto3,enum=api.ValueType" json:"type,omitempty"`
	NumberValue *int32    `protobuf:"varint,2,opt,name=number_value,json=numberValue,proto3,oneof" json:"number_value,omitempty"`
	StringValue *string   `protobuf:"bytes,3,opt,name=string_value,json=stringValue,proto3,oneof" json:"string_value,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_intent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_proto_intent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{0}
}

func (x *Value) GetType() ValueType {
	if x != nil {
		return x.Type
	}
	return ValueType_VALUE_TYPE_UNSPECIFIED
}

func (x *Value) GetNumberValue() int32 {
	if x != nil && x.NumberValue != nil {
		return *x.NumberValue
	}
	return 0
}

func (x *Value) GetStringValue() string {
	if x != nil && x.StringValue != nil {
		return *x.StringValue
	}
	return ""
}

type Intent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   IntentType `protobuf:"varint,1,opt,name=type,proto3,enum=api.IntentType" json:"type,omitempty"`
	Values []*Value   `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *Intent) Reset() {
	*x = Intent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_intent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Intent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Intent) ProtoMessage() {}

func (x *Intent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_intent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Intent.ProtoReflect.Descriptor instead.
func (*Intent) Descriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{1}
}

func (x *Intent) GetType() IntentType {
	if x != nil {
		return x.Type
	}
	return IntentType_INTENT_TYPE_UNSPECIFIED
}

func (x *Intent) GetValues() []*Value {
	if x != nil {
		return x.Values
	}
	return nil
}

type PathRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ipv6DestinationAddress string    `protobuf:"bytes,1,opt,name=ipv6_destination_address,json=ipv6DestinationAddress,proto3" json:"ipv6_destination_address,omitempty"`
	Intents                []*Intent `protobuf:"bytes,2,rep,name=intents,proto3" json:"intents,omitempty"`
}

func (x *PathRequest) Reset() {
	*x = PathRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_intent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathRequest) ProtoMessage() {}

func (x *PathRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_intent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathRequest.ProtoReflect.Descriptor instead.
func (*PathRequest) Descriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{2}
}

func (x *PathRequest) GetIpv6DestinationAddress() string {
	if x != nil {
		return x.Ipv6DestinationAddress
	}
	return ""
}

func (x *PathRequest) GetIntents() []*Intent {
	if x != nil {
		return x.Intents
	}
	return nil
}

type PathResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ipv6DestinationAddress string    `protobuf:"bytes,1,opt,name=ipv6_destination_address,json=ipv6DestinationAddress,proto3" json:"ipv6_destination_address,omitempty"`
	Intents                []*Intent `protobuf:"bytes,2,rep,name=intents,proto3" json:"intents,omitempty"`
	Ipv6SidAddresses       []string  `protobuf:"bytes,3,rep,name=ipv6_sid_addresses,json=ipv6SidAddresses,proto3" json:"ipv6_sid_addresses,omitempty"`
}

func (x *PathResult) Reset() {
	*x = PathResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_intent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PathResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PathResult) ProtoMessage() {}

func (x *PathResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_intent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PathResult.ProtoReflect.Descriptor instead.
func (*PathResult) Descriptor() ([]byte, []int) {
	return file_proto_intent_proto_rawDescGZIP(), []int{3}
}

func (x *PathResult) GetIpv6DestinationAddress() string {
	if x != nil {
		return x.Ipv6DestinationAddress
	}
	return ""
}

func (x *PathResult) GetIntents() []*Intent {
	if x != nil {
		return x.Intents
	}
	return nil
}

func (x *PathResult) GetIpv6SidAddresses() []string {
	if x != nil {
		return x.Ipv6SidAddresses
	}
	return nil
}

var File_proto_intent_proto protoreflect.FileDescriptor

var file_proto_intent_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61, 0x70, 0x69, 0x22, 0x9d, 0x01, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x0c, 0x6e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52,
	0x0b, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x12,
	0x26, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x6e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x51, 0x0a, 0x06, 0x49, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x12, 0x23, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x22, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x6e, 0x0a, 0x0b,
	0x50, 0x61, 0x74, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x38, 0x0a, 0x18, 0x69,
	0x70, 0x76, 0x36, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x69,
	0x70, 0x76, 0x36, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x52, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x9b, 0x01, 0x0a,
	0x0a, 0x50, 0x61, 0x74, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x38, 0x0a, 0x18, 0x69,
	0x70, 0x76, 0x36, 0x5f, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x16, 0x69,
	0x70, 0x76, 0x36, 0x44, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x25, 0x0a, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x52, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2c, 0x0a, 0x12,
	0x69, 0x70, 0x76, 0x36, 0x5f, 0x73, 0x69, 0x64, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x69, 0x70, 0x76, 0x36, 0x53, 0x69,
	0x64, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x2a, 0xc2, 0x01, 0x0a, 0x0a, 0x49,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x17, 0x49, 0x4e, 0x54,
	0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1e, 0x0a, 0x1a, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x49, 0x47, 0x48, 0x5f, 0x42, 0x41, 0x4e, 0x44, 0x57,
	0x49, 0x44, 0x54, 0x48, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x4f, 0x57, 0x5f, 0x42, 0x41, 0x4e, 0x44, 0x57, 0x49,
	0x44, 0x54, 0x48, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x4c, 0x4f, 0x57, 0x5f, 0x4c, 0x41, 0x54, 0x45, 0x4e, 0x43, 0x59,
	0x10, 0x03, 0x12, 0x1f, 0x0a, 0x1b, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4c, 0x4f, 0x57, 0x5f, 0x50, 0x41, 0x43, 0x4b, 0x45, 0x54, 0x5f, 0x4c, 0x4f, 0x53,
	0x53, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4c, 0x4f, 0x57, 0x5f, 0x4a, 0x49, 0x54, 0x54, 0x45, 0x52, 0x10, 0x05, 0x2a,
	0x8c, 0x01, 0x0a, 0x09, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a,
	0x16, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14, 0x56, 0x41, 0x4c,
	0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x5f, 0x56, 0x41, 0x4c, 0x55,
	0x45, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x4d, 0x41, 0x58, 0x5f, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x02, 0x12, 0x12, 0x0a,
	0x0e, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x46, 0x43, 0x10,
	0x03, 0x12, 0x1b, 0x0a, 0x17, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x46, 0x4c, 0x45, 0x58, 0x5f, 0x41, 0x4c, 0x47, 0x4f, 0x5f, 0x4e, 0x52, 0x10, 0x04, 0x32, 0x4a,
	0x0a, 0x10, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x12, 0x36, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x50,
	0x61, 0x74, 0x68, 0x12, 0x10, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x61, 0x74, 0x68, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x50, 0x61, 0x74, 0x68,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x28, 0x01, 0x30, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_intent_proto_rawDescOnce sync.Once
	file_proto_intent_proto_rawDescData = file_proto_intent_proto_rawDesc
)

func file_proto_intent_proto_rawDescGZIP() []byte {
	file_proto_intent_proto_rawDescOnce.Do(func() {
		file_proto_intent_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_intent_proto_rawDescData)
	})
	return file_proto_intent_proto_rawDescData
}

var file_proto_intent_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_intent_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_intent_proto_goTypes = []interface{}{
	(IntentType)(0),     // 0: api.IntentType
	(ValueType)(0),      // 1: api.ValueType
	(*Value)(nil),       // 2: api.Value
	(*Intent)(nil),      // 3: api.Intent
	(*PathRequest)(nil), // 4: api.PathRequest
	(*PathResult)(nil),  // 5: api.PathResult
}
var file_proto_intent_proto_depIdxs = []int32{
	1, // 0: api.Value.type:type_name -> api.ValueType
	0, // 1: api.Intent.type:type_name -> api.IntentType
	2, // 2: api.Intent.values:type_name -> api.Value
	3, // 3: api.PathRequest.intents:type_name -> api.Intent
	3, // 4: api.PathResult.intents:type_name -> api.Intent
	4, // 5: api.IntentController.GetIntentPath:input_type -> api.PathRequest
	5, // 6: api.IntentController.GetIntentPath:output_type -> api.PathResult
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_intent_proto_init() }
func file_proto_intent_proto_init() {
	if File_proto_intent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_intent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_proto_intent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Intent); i {
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
		file_proto_intent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathRequest); i {
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
		file_proto_intent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PathResult); i {
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
	file_proto_intent_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_intent_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_intent_proto_goTypes,
		DependencyIndexes: file_proto_intent_proto_depIdxs,
		EnumInfos:         file_proto_intent_proto_enumTypes,
		MessageInfos:      file_proto_intent_proto_msgTypes,
	}.Build()
	File_proto_intent_proto = out.File
	file_proto_intent_proto_rawDesc = nil
	file_proto_intent_proto_goTypes = nil
	file_proto_intent_proto_depIdxs = nil
}
