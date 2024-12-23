// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.1
// source: protocol/default_users_settings/default_users_settings.proto

package default_users_settings

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UserSettings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// user_uuid uuid рефовода
	UserUuid string            `protobuf:"bytes,1,opt,name=user_uuid,json=userUuid,proto3" json:"user_uuid,omitempty"`
	Settings []*UserPreference `protobuf:"bytes,2,rep,name=settings,proto3" json:"settings,omitempty"`
}

func (x *UserSettings) Reset() {
	*x = UserSettings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserSettings) ProtoMessage() {}

func (x *UserSettings) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserSettings.ProtoReflect.Descriptor instead.
func (*UserSettings) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{0}
}

func (x *UserSettings) GetUserUuid() string {
	if x != nil {
		return x.UserUuid
	}
	return ""
}

func (x *UserSettings) GetSettings() []*UserPreference {
	if x != nil {
		return x.Settings
	}
	return nil
}

type UserPreference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Product string `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
	Coin    string `protobuf:"bytes,2,opt,name=coin,proto3" json:"coin,omitempty"`
	// fee дефолтная комиссия
	Fee float64 `protobuf:"fixed64,3,opt,name=fee,proto3" json:"fee,omitempty"`
	// referral_fee процент от дефолтной комиссии который идёт рефоводу, не может быть больше 100.
	ReferralFee float64 `protobuf:"fixed64,4,opt,name=referral_fee,json=referralFee,proto3" json:"referral_fee,omitempty"`
	// поля ниже нужны только в случае когда отдаём инф. о комиссия
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *UserPreference) Reset() {
	*x = UserPreference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserPreference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserPreference) ProtoMessage() {}

func (x *UserPreference) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserPreference.ProtoReflect.Descriptor instead.
func (*UserPreference) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{1}
}

func (x *UserPreference) GetProduct() string {
	if x != nil {
		return x.Product
	}
	return ""
}

func (x *UserPreference) GetCoin() string {
	if x != nil {
		return x.Coin
	}
	return ""
}

func (x *UserPreference) GetFee() float64 {
	if x != nil {
		return x.Fee
	}
	return 0
}

func (x *UserPreference) GetReferralFee() float64 {
	if x != nil {
		return x.ReferralFee
	}
	return 0
}

func (x *UserPreference) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *UserPreference) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type CreateUsersSettingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*UserSettings `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *CreateUsersSettingsRequest) Reset() {
	*x = CreateUsersSettingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUsersSettingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUsersSettingsRequest) ProtoMessage() {}

func (x *CreateUsersSettingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUsersSettingsRequest.ProtoReflect.Descriptor instead.
func (*CreateUsersSettingsRequest) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUsersSettingsRequest) GetUsers() []*UserSettings {
	if x != nil {
		return x.Users
	}
	return nil
}

type CreateUsersSettingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateUsersSettingsResponse) Reset() {
	*x = CreateUsersSettingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateUsersSettingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUsersSettingsResponse) ProtoMessage() {}

func (x *CreateUsersSettingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUsersSettingsResponse.ProtoReflect.Descriptor instead.
func (*CreateUsersSettingsResponse) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{3}
}

type UpdateUsersSettingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// update_mode:
	// пустое значение - обновляем настройки только у реферала, не трогая его рефчиков
	// all - обновляем настройки только у реферала и у всех его рефчиков учитывая промокоды, если есть промокод то не обновляет.
	// force_all - обновляем настройки только у реферала и у всех его рефчиков, не учитывая промокоды рефчиков.
	UpdateMode string          `protobuf:"bytes,1,opt,name=update_mode,json=updateMode,proto3" json:"update_mode,omitempty"`
	Users      []*UserSettings `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *UpdateUsersSettingsRequest) Reset() {
	*x = UpdateUsersSettingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUsersSettingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUsersSettingsRequest) ProtoMessage() {}

func (x *UpdateUsersSettingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUsersSettingsRequest.ProtoReflect.Descriptor instead.
func (*UpdateUsersSettingsRequest) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateUsersSettingsRequest) GetUpdateMode() string {
	if x != nil {
		return x.UpdateMode
	}
	return ""
}

func (x *UpdateUsersSettingsRequest) GetUsers() []*UserSettings {
	if x != nil {
		return x.Users
	}
	return nil
}

type UpdateUsersSettingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateUsersSettingsResponse) Reset() {
	*x = UpdateUsersSettingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUsersSettingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUsersSettingsResponse) ProtoMessage() {}

func (x *UpdateUsersSettingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUsersSettingsResponse.ProtoReflect.Descriptor instead.
func (*UpdateUsersSettingsResponse) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{5}
}

type GetUsersSettingsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// user_uuids список пользователей, не больше 100 можно запросить
	UserUuids []string `protobuf:"bytes,1,rep,name=user_uuids,json=userUuids,proto3" json:"user_uuids,omitempty"`
	// products список продуктов, если продукты не указано то по всем будет запрос
	Products []string `protobuf:"bytes,2,rep,name=products,proto3" json:"products,omitempty"`
	// coins список монет, если монеты не указаны будет выбора по всем
	Coins []string `protobuf:"bytes,3,rep,name=coins,proto3" json:"coins,omitempty"`
}

func (x *GetUsersSettingsRequest) Reset() {
	*x = GetUsersSettingsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersSettingsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersSettingsRequest) ProtoMessage() {}

func (x *GetUsersSettingsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersSettingsRequest.ProtoReflect.Descriptor instead.
func (*GetUsersSettingsRequest) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{6}
}

func (x *GetUsersSettingsRequest) GetUserUuids() []string {
	if x != nil {
		return x.UserUuids
	}
	return nil
}

func (x *GetUsersSettingsRequest) GetProducts() []string {
	if x != nil {
		return x.Products
	}
	return nil
}

func (x *GetUsersSettingsRequest) GetCoins() []string {
	if x != nil {
		return x.Coins
	}
	return nil
}

type GetUsersSettingsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*UserSettings `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *GetUsersSettingsResponse) Reset() {
	*x = GetUsersSettingsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersSettingsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersSettingsResponse) ProtoMessage() {}

func (x *GetUsersSettingsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protocol_default_users_settings_default_users_settings_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersSettingsResponse.ProtoReflect.Descriptor instead.
func (*GetUsersSettingsResponse) Descriptor() ([]byte, []int) {
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP(), []int{7}
}

func (x *GetUsersSettingsResponse) GetUsers() []*UserSettings {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_protocol_default_users_settings_default_users_settings_proto protoreflect.FileDescriptor

var file_protocol_default_users_settings_default_users_settings_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x73, 0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f,
	0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16,
	0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x55, 0x75, 0x69, 0x64, 0x12, 0x42, 0x0a, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x08,
	0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x22, 0xe9, 0x01, 0x0a, 0x0e, 0x55, 0x73, 0x65,
	0x72, 0x50, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x69, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x65, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x66, 0x65, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x61, 0x6c, 0x5f, 0x66, 0x65, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x72, 0x65, 0x66, 0x65, 0x72, 0x72, 0x61, 0x6c, 0x46, 0x65, 0x65, 0x12, 0x39,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x22, 0x58, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x3a, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x24, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72,
	0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x1d,
	0x0a, 0x1b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x79, 0x0a,
	0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x3a, 0x0a, 0x05,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x64, 0x65,
	0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x73, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x22, 0x1d, 0x0a, 0x1b, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x6a, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x63, 0x6f, 0x69, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x63, 0x6f,
	0x69, 0x6e, 0x73, 0x22, 0x56, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3a, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x32, 0x9e, 0x03, 0x0a, 0x1b,
	0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x80, 0x01, 0x0a, 0x13,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x12, 0x32, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x80,
	0x01, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x32, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x64, 0x65, 0x66,
	0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x79, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74,
	0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x2f, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x30, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x60, 0x5a, 0x5e,
	0x63, 0x6f, 0x64, 0x65, 0x2e, 0x65, 0x6d, 0x63, 0x64, 0x74, 0x65, 0x63, 0x68, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x65, 0x6d, 0x63, 0x64, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x72,
	0x65, 0x66, 0x65, 0x72, 0x72, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x2f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73,
	0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x5f, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protocol_default_users_settings_default_users_settings_proto_rawDescOnce sync.Once
	file_protocol_default_users_settings_default_users_settings_proto_rawDescData = file_protocol_default_users_settings_default_users_settings_proto_rawDesc
)

func file_protocol_default_users_settings_default_users_settings_proto_rawDescGZIP() []byte {
	file_protocol_default_users_settings_default_users_settings_proto_rawDescOnce.Do(func() {
		file_protocol_default_users_settings_default_users_settings_proto_rawDescData = protoimpl.X.CompressGZIP(file_protocol_default_users_settings_default_users_settings_proto_rawDescData)
	})
	return file_protocol_default_users_settings_default_users_settings_proto_rawDescData
}

var file_protocol_default_users_settings_default_users_settings_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_protocol_default_users_settings_default_users_settings_proto_goTypes = []any{
	(*UserSettings)(nil),                // 0: default_users_settings.UserSettings
	(*UserPreference)(nil),              // 1: default_users_settings.UserPreference
	(*CreateUsersSettingsRequest)(nil),  // 2: default_users_settings.CreateUsersSettingsRequest
	(*CreateUsersSettingsResponse)(nil), // 3: default_users_settings.CreateUsersSettingsResponse
	(*UpdateUsersSettingsRequest)(nil),  // 4: default_users_settings.UpdateUsersSettingsRequest
	(*UpdateUsersSettingsResponse)(nil), // 5: default_users_settings.UpdateUsersSettingsResponse
	(*GetUsersSettingsRequest)(nil),     // 6: default_users_settings.GetUsersSettingsRequest
	(*GetUsersSettingsResponse)(nil),    // 7: default_users_settings.GetUsersSettingsResponse
	(*timestamppb.Timestamp)(nil),       // 8: google.protobuf.Timestamp
}
var file_protocol_default_users_settings_default_users_settings_proto_depIdxs = []int32{
	1, // 0: default_users_settings.UserSettings.settings:type_name -> default_users_settings.UserPreference
	8, // 1: default_users_settings.UserPreference.created_at:type_name -> google.protobuf.Timestamp
	8, // 2: default_users_settings.UserPreference.updated_at:type_name -> google.protobuf.Timestamp
	0, // 3: default_users_settings.CreateUsersSettingsRequest.users:type_name -> default_users_settings.UserSettings
	0, // 4: default_users_settings.UpdateUsersSettingsRequest.users:type_name -> default_users_settings.UserSettings
	0, // 5: default_users_settings.GetUsersSettingsResponse.users:type_name -> default_users_settings.UserSettings
	2, // 6: default_users_settings.DefaultUsersSettingsService.CreateUsersSettings:input_type -> default_users_settings.CreateUsersSettingsRequest
	4, // 7: default_users_settings.DefaultUsersSettingsService.UpdateUsersSettings:input_type -> default_users_settings.UpdateUsersSettingsRequest
	6, // 8: default_users_settings.DefaultUsersSettingsService.GetUsersSettings:input_type -> default_users_settings.GetUsersSettingsRequest
	3, // 9: default_users_settings.DefaultUsersSettingsService.CreateUsersSettings:output_type -> default_users_settings.CreateUsersSettingsResponse
	5, // 10: default_users_settings.DefaultUsersSettingsService.UpdateUsersSettings:output_type -> default_users_settings.UpdateUsersSettingsResponse
	7, // 11: default_users_settings.DefaultUsersSettingsService.GetUsersSettings:output_type -> default_users_settings.GetUsersSettingsResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_protocol_default_users_settings_default_users_settings_proto_init() }
func file_protocol_default_users_settings_default_users_settings_proto_init() {
	if File_protocol_default_users_settings_default_users_settings_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*UserSettings); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UserPreference); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*CreateUsersSettingsRequest); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*CreateUsersSettingsResponse); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateUsersSettingsRequest); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateUsersSettingsResponse); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*GetUsersSettingsRequest); i {
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
		file_protocol_default_users_settings_default_users_settings_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*GetUsersSettingsResponse); i {
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
			RawDescriptor: file_protocol_default_users_settings_default_users_settings_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protocol_default_users_settings_default_users_settings_proto_goTypes,
		DependencyIndexes: file_protocol_default_users_settings_default_users_settings_proto_depIdxs,
		MessageInfos:      file_protocol_default_users_settings_default_users_settings_proto_msgTypes,
	}.Build()
	File_protocol_default_users_settings_default_users_settings_proto = out.File
	file_protocol_default_users_settings_default_users_settings_proto_rawDesc = nil
	file_protocol_default_users_settings_default_users_settings_proto_goTypes = nil
	file_protocol_default_users_settings_default_users_settings_proto_depIdxs = nil
}
