// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messageid.proto

package PB

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message int32

const (
	Message_None  Message = 0
	Message_Multi Message = 1
	// 消息(系统/世界聊天)部分
	Message_Chat_ChatRequest          Message = 20403
	Message_Chat_ChatTachChatGoup     Message = 20404
	Message_Chat_SetChatChannelSwitch Message = 20406
)

var Message_name = map[int32]string{
	0:     "None",
	1:     "Multi",
	20403: "Chat_ChatRequest",
	20404: "Chat_ChatTachChatGoup",
	20406: "Chat_SetChatChannelSwitch",
}
var Message_value = map[string]int32{
	"None":                      0,
	"Multi":                     1,
	"Chat_ChatRequest":          20403,
	"Chat_ChatTachChatGoup":     20404,
	"Chat_SetChatChannelSwitch": 20406,
}

func (x Message) String() string {
	return proto.EnumName(Message_name, int32(x))
}
func (Message) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_messageid_1dffe8107f6004bc, []int{0}
}

func init() {
	proto.RegisterEnum("PB.Message", Message_name, Message_value)
}

func init() { proto.RegisterFile("messageid.proto", fileDescriptor_messageid_1dffe8107f6004bc) }

var fileDescriptor_messageid_1dffe8107f6004bc = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcf, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xcd, 0x4c, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x0a, 0x70, 0xd2,
	0x2a, 0xe1, 0x62, 0xf7, 0x85, 0x08, 0x0b, 0x71, 0x70, 0xb1, 0xf8, 0xe5, 0xe7, 0xa5, 0x0a, 0x30,
	0x08, 0x71, 0x72, 0xb1, 0xfa, 0x96, 0xe6, 0x94, 0x64, 0x0a, 0x30, 0x0a, 0x89, 0x71, 0x09, 0x38,
	0x67, 0x24, 0x96, 0xc4, 0x83, 0x88, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x81, 0xcd, 0xf3,
	0x19, 0x85, 0xa4, 0xb9, 0x44, 0xe1, 0xe2, 0x21, 0x89, 0xc9, 0x19, 0x20, 0xda, 0x3d, 0xbf, 0xb4,
	0x40, 0x60, 0xcb, 0x7c, 0x46, 0x21, 0x79, 0x2e, 0x49, 0xb0, 0x64, 0x70, 0x6a, 0x09, 0x88, 0x76,
	0xce, 0x48, 0xcc, 0xcb, 0x4b, 0xcd, 0x09, 0x2e, 0xcf, 0x2c, 0x49, 0xce, 0x10, 0xd8, 0x36, 0x9f,
	0x31, 0x89, 0x0d, 0xec, 0x00, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfa, 0x9b, 0xbd, 0xea,
	0x93, 0x00, 0x00, 0x00,
}
