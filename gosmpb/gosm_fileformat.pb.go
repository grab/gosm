// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gosm_fileformat.proto

package gosmpb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Blob struct {
	Raw     []byte `protobuf:"bytes,1,opt,name=raw" json:"raw,omitempty"`
	RawSize *int32 `protobuf:"varint,2,opt,name=raw_size,json=rawSize" json:"raw_size,omitempty"`
	// Possible compressed versions of the data.
	ZlibData []byte `protobuf:"bytes,3,opt,name=zlib_data,json=zlibData" json:"zlib_data,omitempty"`
	// PROPOSED feature for LZMA compressed data. SUPPORT IS NOT REQUIRED.
	LzmaData []byte `protobuf:"bytes,4,opt,name=lzma_data,json=lzmaData" json:"lzma_data,omitempty"`
	// Formerly used for bzip2 compressed data. Depreciated in 2010.
	OBSOLETEBzip2Data    []byte   `protobuf:"bytes,5,opt,name=OBSOLETE_bzip2_data,json=OBSOLETEBzip2Data" json:"OBSOLETE_bzip2_data,omitempty"` // Deprecated: Do not use.
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Blob) Reset()         { *m = Blob{} }
func (m *Blob) String() string { return proto.CompactTextString(m) }
func (*Blob) ProtoMessage()    {}
func (*Blob) Descriptor() ([]byte, []int) {
	return fileDescriptor_2462ba141d092c4b, []int{0}
}

func (m *Blob) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Blob.Unmarshal(m, b)
}
func (m *Blob) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Blob.Marshal(b, m, deterministic)
}
func (m *Blob) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Blob.Merge(m, src)
}
func (m *Blob) XXX_Size() int {
	return xxx_messageInfo_Blob.Size(m)
}
func (m *Blob) XXX_DiscardUnknown() {
	xxx_messageInfo_Blob.DiscardUnknown(m)
}

var xxx_messageInfo_Blob proto.InternalMessageInfo

func (m *Blob) GetRaw() []byte {
	if m != nil {
		return m.Raw
	}
	return nil
}

func (m *Blob) GetRawSize() int32 {
	if m != nil && m.RawSize != nil {
		return *m.RawSize
	}
	return 0
}

func (m *Blob) GetZlibData() []byte {
	if m != nil {
		return m.ZlibData
	}
	return nil
}

func (m *Blob) GetLzmaData() []byte {
	if m != nil {
		return m.LzmaData
	}
	return nil
}

// Deprecated: Do not use.
func (m *Blob) GetOBSOLETEBzip2Data() []byte {
	if m != nil {
		return m.OBSOLETEBzip2Data
	}
	return nil
}

type BlobHeader struct {
	Type                 *string  `protobuf:"bytes,1,req,name=type" json:"type,omitempty"`
	Indexdata            []byte   `protobuf:"bytes,2,opt,name=indexdata" json:"indexdata,omitempty"`
	Datasize             *int32   `protobuf:"varint,3,req,name=datasize" json:"datasize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlobHeader) Reset()         { *m = BlobHeader{} }
func (m *BlobHeader) String() string { return proto.CompactTextString(m) }
func (*BlobHeader) ProtoMessage()    {}
func (*BlobHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_2462ba141d092c4b, []int{1}
}

func (m *BlobHeader) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlobHeader.Unmarshal(m, b)
}
func (m *BlobHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlobHeader.Marshal(b, m, deterministic)
}
func (m *BlobHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlobHeader.Merge(m, src)
}
func (m *BlobHeader) XXX_Size() int {
	return xxx_messageInfo_BlobHeader.Size(m)
}
func (m *BlobHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlobHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlobHeader proto.InternalMessageInfo

func (m *BlobHeader) GetType() string {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return ""
}

func (m *BlobHeader) GetIndexdata() []byte {
	if m != nil {
		return m.Indexdata
	}
	return nil
}

func (m *BlobHeader) GetDatasize() int32 {
	if m != nil && m.Datasize != nil {
		return *m.Datasize
	}
	return 0
}

func init() {
	proto.RegisterType((*Blob)(nil), "gosmpb.Blob")
	proto.RegisterType((*BlobHeader)(nil), "gosmpb.BlobHeader")
}

func init() { proto.RegisterFile("gosm_fileformat.proto", fileDescriptor_2462ba141d092c4b) }

var fileDescriptor_2462ba141d092c4b = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcd, 0x4a, 0xc4, 0x30,
	0x10, 0xc7, 0x49, 0xda, 0x6a, 0x3b, 0x78, 0xd0, 0x88, 0x10, 0x3f, 0x0e, 0x65, 0x4f, 0x3d, 0x79,
	0xd8, 0x47, 0x08, 0x2e, 0x78, 0x10, 0x16, 0xb2, 0x9e, 0xf6, 0x52, 0xa6, 0x34, 0x2b, 0x81, 0xd6,
	0x94, 0x34, 0x50, 0xcd, 0xdb, 0xf8, 0xa6, 0x32, 0xa9, 0x1f, 0xa7, 0x4c, 0x7e, 0xbf, 0xf9, 0xe7,
	0x63, 0xe0, 0xe6, 0xcd, 0xcd, 0x63, 0x7b, 0xb2, 0x83, 0x39, 0x39, 0x3f, 0x62, 0x78, 0x9c, 0xbc,
	0x0b, 0x4e, 0x9c, 0x11, 0x9e, 0xba, 0xcd, 0x17, 0x83, 0x5c, 0x0d, 0xae, 0x13, 0x97, 0x90, 0x79,
	0x5c, 0x24, 0xab, 0x59, 0x73, 0xa1, 0xa9, 0x14, 0xb7, 0x50, 0x7a, 0x5c, 0xda, 0xd9, 0x46, 0x23,
	0x79, 0xcd, 0x9a, 0x42, 0x9f, 0x7b, 0x5c, 0x0e, 0x36, 0x1a, 0x71, 0x0f, 0x55, 0x1c, 0x6c, 0xd7,
	0xf6, 0x18, 0x50, 0x66, 0x29, 0x52, 0x12, 0x78, 0xc2, 0x80, 0x24, 0x87, 0x38, 0xe2, 0x2a, 0xf3,
	0x55, 0x12, 0x48, 0x72, 0x0b, 0xd7, 0x7b, 0x75, 0xd8, 0xbf, 0xec, 0x5e, 0x77, 0x6d, 0x17, 0xed,
	0xb4, 0x5d, 0xdb, 0x0a, 0x6a, 0x53, 0x5c, 0x32, 0x7d, 0xf5, 0xab, 0x15, 0x59, 0xca, 0x6c, 0x8e,
	0x00, 0xf4, 0xc4, 0x67, 0x83, 0xbd, 0xf1, 0x42, 0x40, 0x1e, 0x3e, 0x27, 0x23, 0x59, 0xcd, 0x9b,
	0x4a, 0xa7, 0x5a, 0x3c, 0x40, 0x65, 0xdf, 0x7b, 0xf3, 0x91, 0xce, 0xe2, 0xe9, 0xca, 0x7f, 0x20,
	0xee, 0xa0, 0xa4, 0x35, 0x7d, 0x24, 0xab, 0x79, 0x53, 0xe8, 0xbf, 0xbd, 0x2a, 0x8f, 0x3f, 0x93,
	0xf8, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x87, 0x17, 0x7e, 0x29, 0x01, 0x00, 0x00,
}
