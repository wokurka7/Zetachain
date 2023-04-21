// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: zetacore/crosschain/node_account.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	common "github.com/zeta-chain/zetacore/common"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type NodeStatus int32

const (
	NodeStatus_Unknown     NodeStatus = 0
	NodeStatus_Whitelisted NodeStatus = 1
	NodeStatus_Standby     NodeStatus = 2
	NodeStatus_Ready       NodeStatus = 3
	NodeStatus_Active      NodeStatus = 4
	NodeStatus_Disabled    NodeStatus = 5
)

// Genesis
// NonGenesis

var NodeStatus_name = map[int32]string{
	0: "Unknown",
	1: "Whitelisted",
	2: "Standby",
	3: "Ready",
	4: "Active",
	5: "Disabled",
}

var NodeStatus_value = map[string]int32{
	"Unknown":     0,
	"Whitelisted": 1,
	"Standby":     2,
	"Ready":       3,
	"Active":      4,
	"Disabled":    5,
}

func (x NodeStatus) String() string {
	return proto.EnumName(NodeStatus_name, int32(x))
}

func (NodeStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_10b3c7d9388b75aa, []int{0}
}

type NodeAccount struct {
	Creator          string            `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	TssSignerAddress string            `protobuf:"bytes,2,opt,name=tssSignerAddress,proto3" json:"tssSignerAddress,omitempty"`
	PubkeySet        *common.PubKeySet `protobuf:"bytes,3,opt,name=pubkeySet,proto3" json:"pubkeySet,omitempty"`
	NodeStatus       NodeStatus        `protobuf:"varint,4,opt,name=nodeStatus,proto3,enum=zetachain.zetacore.crosschain.NodeStatus" json:"nodeStatus,omitempty"`
}

func (m *NodeAccount) Reset()         { *m = NodeAccount{} }
func (m *NodeAccount) String() string { return proto.CompactTextString(m) }
func (*NodeAccount) ProtoMessage()    {}
func (*NodeAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_10b3c7d9388b75aa, []int{0}
}
func (m *NodeAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NodeAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NodeAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NodeAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeAccount.Merge(m, src)
}
func (m *NodeAccount) XXX_Size() int {
	return m.Size()
}
func (m *NodeAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeAccount.DiscardUnknown(m)
}

var xxx_messageInfo_NodeAccount proto.InternalMessageInfo

func (m *NodeAccount) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *NodeAccount) GetTssSignerAddress() string {
	if m != nil {
		return m.TssSignerAddress
	}
	return ""
}

func (m *NodeAccount) GetPubkeySet() *common.PubKeySet {
	if m != nil {
		return m.PubkeySet
	}
	return nil
}

func (m *NodeAccount) GetNodeStatus() NodeStatus {
	if m != nil {
		return m.NodeStatus
	}
	return NodeStatus_Unknown
}

func init() {
	proto.RegisterEnum("zetacore.crosschain.NodeStatus", NodeStatus_name, NodeStatus_value)
	proto.RegisterType((*NodeAccount)(nil), "zetacore.crosschain.NodeAccount")
}

func init() {
	proto.RegisterFile("zetacore/crosschain/node_account.proto", fileDescriptor_10b3c7d9388b75aa)
}

var fileDescriptor_ea30ee4c0fac150c = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xcd, 0x6a, 0xdb, 0x40,
	0x18, 0xd4, 0xfa, 0xb7, 0xfe, 0x54, 0x5a, 0x75, 0xdb, 0x83, 0x30, 0x58, 0x98, 0x9e, 0x5c, 0x43,
	0x25, 0xea, 0x3e, 0x81, 0x43, 0x2e, 0xc1, 0x60, 0x82, 0x44, 0x08, 0xe4, 0x12, 0x56, 0xda, 0x0f,
	0x59, 0xd8, 0xde, 0x35, 0xda, 0x55, 0x12, 0xe5, 0x29, 0xf2, 0x10, 0x39, 0xe4, 0x51, 0x72, 0x8b,
	0x8f, 0x39, 0x06, 0xfb, 0x45, 0x82, 0xa4, 0x38, 0x36, 0x04, 0x72, 0xda, 0x6f, 0xbf, 0x99, 0xd9,
	0x9d, 0x61, 0xa0, 0x17, 0xa5, 0x52, 0xa9, 0x68, 0xc6, 0x12, 0xe1, 0x09, 0xc9, 0xf1, 0x92, 0x45,
	0x91, 0xcc, 0x84, 0x76, 0x57, 0xa9, 0xd4, 0x92, 0xf6, 0x6e, 0x51, 0xb3, 0x12, 0x75, 0xcb, 0x49,
	0xa6, 0xe8, 0xee, 0x15, 0xdd, 0x5f, 0xb1, 0x8c, 0x65, 0xc9, 0xf4, 0x8a, 0xa9, 0x12, 0x75, 0x7f,
	0x46, 0x72, 0xb9, 0x94, 0xc2, 0xab, 0x8e, 0x6a, 0xf9, 0xfb, 0x89, 0x80, 0x39, 0x95, 0x1c, 0xc7,
	0xd5, 0xfb, 0xd4, 0x86, 0x76, 0x94, 0x22, 0xd3, 0x32, 0xb5, 0x49, 0x9f, 0x0c, 0x3a, 0xfe, 0xee,
	0x4a, 0x87, 0x60, 0x69, 0xa5, 0x82, 0x24, 0x16, 0x98, 0x8e, 0x39, 0x4f, 0x51, 0x29, 0xbb, 0x56,
	0x52, 0x3e, 0xec, 0xa9, 0x07, 0x9d, 0x55, 0x16, 0xce, 0x31, 0x0f, 0x50, 0xdb, 0xf5, 0x3e, 0x19,
	0x98, 0xa3, 0x1f, 0xee, 0xdb, 0xbf, 0xa7, 0x59, 0x38, 0x29, 0x01, 0x7f, 0xcf, 0xa1, 0x27, 0x00,
	0x45, 0xcc, 0x40, 0x33, 0x9d, 0x29, 0xbb, 0xd1, 0x27, 0x83, 0x6f, 0xa3, 0x3f, 0xee, 0xa7, 0x29,
	0xdd, 0xe9, 0xbb, 0xc0, 0x3f, 0x10, 0x0f, 0x43, 0x80, 0x3d, 0x42, 0x4d, 0x68, 0x9f, 0x89, 0xb9,
	0x90, 0xd7, 0xc2, 0x32, 0xe8, 0x77, 0x30, 0xcf, 0x67, 0x89, 0xc6, 0x45, 0xa2, 0x34, 0x72, 0x8b,
	0x14, 0x68, 0xa0, 0x99, 0xe0, 0x61, 0x6e, 0xd5, 0x68, 0x07, 0x9a, 0x3e, 0x32, 0x9e, 0x5b, 0x75,
	0x0a, 0xd0, 0x1a, 0x47, 0x3a, 0xb9, 0x42, 0xab, 0x41, 0xbf, 0xc2, 0x97, 0xe3, 0x44, 0xb1, 0x70,
	0x81, 0xdc, 0x6a, 0x76, 0x1b, 0x0f, 0xf7, 0x0e, 0x39, 0x9a, 0x3c, 0x6e, 0x1c, 0xb2, 0xde, 0x38,
	0xe4, 0x65, 0xe3, 0x90, 0xbb, 0xad, 0x63, 0xac, 0xb7, 0x8e, 0xf1, 0xbc, 0x75, 0x8c, 0x8b, 0x7f,
	0x71, 0xa2, 0x67, 0x59, 0x58, 0x84, 0xf5, 0x0a, 0xd3, 0x7f, 0xab, 0x0e, 0x77, 0xfe, 0xbd, 0x1b,
	0xef, 0xa0, 0x59, 0x9d, 0xaf, 0x50, 0x85, 0xad, 0xb2, 0x89, 0xff, 0xaf, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x98, 0x54, 0xb2, 0xff, 0xf4, 0x01, 0x00, 0x00,
}

func (m *NodeAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NodeAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NodeAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NodeStatus != 0 {
		i = encodeVarintNodeAccount(dAtA, i, uint64(m.NodeStatus))
		i--
		dAtA[i] = 0x20
	}
	if m.PubkeySet != nil {
		{
			size, err := m.PubkeySet.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintNodeAccount(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TssSignerAddress) > 0 {
		i -= len(m.TssSignerAddress)
		copy(dAtA[i:], m.TssSignerAddress)
		i = encodeVarintNodeAccount(dAtA, i, uint64(len(m.TssSignerAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintNodeAccount(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintNodeAccount(dAtA []byte, offset int, v uint64) int {
	offset -= sovNodeAccount(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *NodeAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	l = len(m.TssSignerAddress)
	if l > 0 {
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	if m.PubkeySet != nil {
		l = m.PubkeySet.Size()
		n += 1 + l + sovNodeAccount(uint64(l))
	}
	if m.NodeStatus != 0 {
		n += 1 + sovNodeAccount(uint64(m.NodeStatus))
	}
	return n
}

func sovNodeAccount(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozNodeAccount(x uint64) (n int) {
	return sovNodeAccount(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *NodeAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowNodeAccount
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: NodeAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NodeAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TssSignerAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TssSignerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubkeySet", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthNodeAccount
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubkeySet == nil {
				m.PubkeySet = &common.PubKeySet{}
			}
			if err := m.PubkeySet.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeStatus", wireType)
			}
			m.NodeStatus = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NodeStatus |= NodeStatus(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipNodeAccount(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthNodeAccount
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipNodeAccount(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowNodeAccount
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowNodeAccount
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthNodeAccount
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupNodeAccount
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthNodeAccount
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthNodeAccount        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowNodeAccount          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupNodeAccount = fmt.Errorf("proto: unexpected end of group")
)
