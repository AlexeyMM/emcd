// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.29.0
// source: admin.proto

package adminpb

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type MerchantTariff struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UpperFee string `protobuf:"bytes,1,opt,name=upper_fee,json=upperFee,proto3" json:"upper_fee,omitempty"` // positive decimal
	LowerFee string `protobuf:"bytes,2,opt,name=lower_fee,json=lowerFee,proto3" json:"lower_fee,omitempty"` // positive decimal
	MinPay   string `protobuf:"bytes,3,opt,name=min_pay,json=minPay,proto3" json:"min_pay,omitempty"`       // positive decimal
	MaxPay   string `protobuf:"bytes,4,opt,name=max_pay,json=maxPay,proto3" json:"max_pay,omitempty"`       // positive decimal
}

func (x *MerchantTariff) Reset() {
	*x = MerchantTariff{}
	mi := &file_admin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MerchantTariff) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerchantTariff) ProtoMessage() {}

func (x *MerchantTariff) ProtoReflect() protoreflect.Message {
	mi := &file_admin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerchantTariff.ProtoReflect.Descriptor instead.
func (*MerchantTariff) Descriptor() ([]byte, []int) {
	return file_admin_proto_rawDescGZIP(), []int{0}
}

func (x *MerchantTariff) GetUpperFee() string {
	if x != nil {
		return x.UpperFee
	}
	return ""
}

func (x *MerchantTariff) GetLowerFee() string {
	if x != nil {
		return x.LowerFee
	}
	return ""
}

func (x *MerchantTariff) GetMinPay() string {
	if x != nil {
		return x.MinPay
	}
	return ""
}

func (x *MerchantTariff) GetMaxPay() string {
	if x != nil {
		return x.MaxPay
	}
	return ""
}

type CreateMerchantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string          `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Tariff *MerchantTariff `protobuf:"bytes,2,opt,name=tariff,proto3" json:"tariff,omitempty"`
}

func (x *CreateMerchantRequest) Reset() {
	*x = CreateMerchantRequest{}
	mi := &file_admin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMerchantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMerchantRequest) ProtoMessage() {}

func (x *CreateMerchantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_admin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMerchantRequest.ProtoReflect.Descriptor instead.
func (*CreateMerchantRequest) Descriptor() ([]byte, []int) {
	return file_admin_proto_rawDescGZIP(), []int{1}
}

func (x *CreateMerchantRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CreateMerchantRequest) GetTariff() *MerchantTariff {
	if x != nil {
		return x.Tariff
	}
	return nil
}

type CreateMerchantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateMerchantResponse) Reset() {
	*x = CreateMerchantResponse{}
	mi := &file_admin_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateMerchantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMerchantResponse) ProtoMessage() {}

func (x *CreateMerchantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_admin_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMerchantResponse.ProtoReflect.Descriptor instead.
func (*CreateMerchantResponse) Descriptor() ([]byte, []int) {
	return file_admin_proto_rawDescGZIP(), []int{2}
}

var File_admin_proto protoreflect.FileDescriptor

var file_admin_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x1a, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x33,
	0x72, 0x64, 0x2d, 0x70, 0x61, 0x72, 0x74, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf8, 0x03, 0x0a, 0x0e, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x54, 0x61, 0x72, 0x69, 0x66, 0x66, 0x12, 0x7c, 0x0a, 0x09, 0x75, 0x70, 0x70, 0x65,
	0x72, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x5f, 0xba, 0x48, 0x5c,
	0xba, 0x01, 0x59, 0x0a, 0x11, 0x75, 0x70, 0x70, 0x65, 0x72, 0x5f, 0x66, 0x65, 0x65, 0x2e, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x12, 0x22, 0x55, 0x70, 0x70, 0x65, 0x72, 0x20, 0x66, 0x65,
	0x65, 0x20, 0x6d, 0x75, 0x73, 0x74, 0x20, 0x62, 0x65, 0x20, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x76, 0x65, 0x20, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x1a, 0x20, 0x74, 0x68, 0x69, 0x73,
	0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5e, 0x5c, 0x5c, 0x64, 0x2b, 0x28,
	0x5c, 0x5c, 0x2e, 0x5c, 0x5c, 0x64, 0x2b, 0x29, 0x3f, 0x24, 0x27, 0x29, 0x52, 0x08, 0x75, 0x70,
	0x70, 0x65, 0x72, 0x46, 0x65, 0x65, 0x12, 0x7c, 0x0a, 0x09, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x5f,
	0x66, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x5f, 0xba, 0x48, 0x5c, 0xba, 0x01,
	0x59, 0x0a, 0x11, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x6c, 0x6f, 0x77, 0x65, 0x72,
	0x5f, 0x66, 0x65, 0x65, 0x12, 0x22, 0x4c, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x66, 0x65, 0x65, 0x20,
	0x6d, 0x75, 0x73, 0x74, 0x20, 0x62, 0x65, 0x20, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65,
	0x20, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x1a, 0x20, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5e, 0x5c, 0x5c, 0x64, 0x2b, 0x28, 0x5c, 0x5c,
	0x2e, 0x5c, 0x5c, 0x64, 0x2b, 0x29, 0x3f, 0x24, 0x27, 0x29, 0x52, 0x08, 0x6c, 0x6f, 0x77, 0x65,
	0x72, 0x46, 0x65, 0x65, 0x12, 0x74, 0x0a, 0x07, 0x6d, 0x69, 0x6e, 0x5f, 0x70, 0x61, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x5b, 0xba, 0x48, 0x58, 0xba, 0x01, 0x55, 0x0a, 0x0f, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x6d, 0x69, 0x6e, 0x5f, 0x70, 0x61, 0x79, 0x12, 0x20,
	0x4d, 0x69, 0x6e, 0x20, 0x70, 0x61, 0x79, 0x20, 0x6d, 0x75, 0x73, 0x74, 0x20, 0x62, 0x65, 0x20,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x20, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c,
	0x1a, 0x20, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x73, 0x28, 0x27,
	0x5e, 0x5c, 0x5c, 0x64, 0x2b, 0x28, 0x5c, 0x5c, 0x2e, 0x5c, 0x5c, 0x64, 0x2b, 0x29, 0x3f, 0x24,
	0x27, 0x29, 0x52, 0x06, 0x6d, 0x69, 0x6e, 0x50, 0x61, 0x79, 0x12, 0x74, 0x0a, 0x07, 0x6d, 0x61,
	0x78, 0x5f, 0x70, 0x61, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x5b, 0xba, 0x48, 0x58,
	0xba, 0x01, 0x55, 0x0a, 0x0f, 0x64, 0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x2e, 0x6d, 0x61, 0x78,
	0x5f, 0x70, 0x61, 0x79, 0x12, 0x20, 0x4d, 0x61, 0x78, 0x20, 0x70, 0x61, 0x79, 0x20, 0x6d, 0x75,
	0x73, 0x74, 0x20, 0x62, 0x65, 0x20, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x20, 0x64,
	0x65, 0x63, 0x69, 0x6d, 0x61, 0x6c, 0x1a, 0x20, 0x74, 0x68, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x65, 0x73, 0x28, 0x27, 0x5e, 0x5c, 0x5c, 0x64, 0x2b, 0x28, 0x5c, 0x5c, 0x2e, 0x5c,
	0x5c, 0x64, 0x2b, 0x29, 0x3f, 0x24, 0x27, 0x29, 0x52, 0x06, 0x6d, 0x61, 0x78, 0x50, 0x61, 0x79,
	0x22, 0x71, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xba, 0x48, 0x05, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x06,
	0x74, 0x61, 0x72, 0x69, 0x66, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x54, 0x61, 0x72,
	0x69, 0x66, 0x66, 0x42, 0x06, 0xba, 0x48, 0x03, 0xc8, 0x01, 0x01, 0x52, 0x06, 0x74, 0x61, 0x72,
	0x69, 0x66, 0x66, 0x22, 0x18, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x72,
	0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x65, 0x0a,
	0x14, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d,
	0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x72, 0x63, 0x68, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x33, 0x5a, 0x31, 0x63, 0x6f, 0x64, 0x65, 0x2e, 0x65, 0x6d, 0x63,
	0x64, 0x74, 0x65, 0x63, 0x68, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x32, 0x62, 0x2f, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_admin_proto_rawDescOnce sync.Once
	file_admin_proto_rawDescData = file_admin_proto_rawDesc
)

func file_admin_proto_rawDescGZIP() []byte {
	file_admin_proto_rawDescOnce.Do(func() {
		file_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_admin_proto_rawDescData)
	})
	return file_admin_proto_rawDescData
}

var file_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_admin_proto_goTypes = []any{
	(*MerchantTariff)(nil),         // 0: admin.MerchantTariff
	(*CreateMerchantRequest)(nil),  // 1: admin.CreateMerchantRequest
	(*CreateMerchantResponse)(nil), // 2: admin.CreateMerchantResponse
}
var file_admin_proto_depIdxs = []int32{
	0, // 0: admin.CreateMerchantRequest.tariff:type_name -> admin.MerchantTariff
	1, // 1: admin.MerchantAdminService.CreateMerchant:input_type -> admin.CreateMerchantRequest
	2, // 2: admin.MerchantAdminService.CreateMerchant:output_type -> admin.CreateMerchantResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_admin_proto_init() }
func file_admin_proto_init() {
	if File_admin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_admin_proto_goTypes,
		DependencyIndexes: file_admin_proto_depIdxs,
		MessageInfos:      file_admin_proto_msgTypes,
	}.Build()
	File_admin_proto = out.File
	file_admin_proto_rawDesc = nil
	file_admin_proto_goTypes = nil
	file_admin_proto_depIdxs = nil
}
