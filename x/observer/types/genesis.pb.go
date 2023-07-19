// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: observer/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type GenesisState struct {
	Ballots         []*Ballot         `protobuf:"bytes,1,rep,name=ballots,proto3" json:"ballots,omitempty"`
	Observers       []*ObserverMapper `protobuf:"bytes,2,rep,name=observers,proto3" json:"observers,omitempty"`
	NodeAccountList []*NodeAccount    `protobuf:"bytes,3,rep,name=nodeAccountList,proto3" json:"nodeAccountList,omitempty"`
	PermissionFlags *PermissionFlags  `protobuf:"bytes,4,opt,name=permissionFlags,proto3" json:"permissionFlags,omitempty"`
	Params          *Params           `protobuf:"bytes,5,opt,name=params,proto3" json:"params,omitempty"`
	Keygen          *Keygen           `protobuf:"bytes,6,opt,name=keygen,proto3" json:"keygen,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_15ea8c9d44da7399, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetBallots() []*Ballot {
	if m != nil {
		return m.Ballots
	}
	return nil
}

func (m *GenesisState) GetObservers() []*ObserverMapper {
	if m != nil {
		return m.Observers
	}
	return nil
}

func (m *GenesisState) GetNodeAccountList() []*NodeAccount {
	if m != nil {
		return m.NodeAccountList
	}
	return nil
}

func (m *GenesisState) GetPermissionFlags() *PermissionFlags {
	if m != nil {
		return m.PermissionFlags
	}
	return nil
}

func (m *GenesisState) GetParams() *Params {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *GenesisState) GetKeygen() *Keygen {
	if m != nil {
		return m.Keygen
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "zetachain.zetacore.observer.GenesisState")
}

func init() { proto.RegisterFile("observer/genesis.proto", fileDescriptor_15ea8c9d44da7399) }

var fileDescriptor_15ea8c9d44da7399 = []byte{
	// 364 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xcb, 0x4a, 0xc3, 0x40,
	0x14, 0x86, 0x1b, 0xab, 0x15, 0xa7, 0x42, 0x21, 0x78, 0x09, 0x2d, 0xc4, 0xa2, 0x9b, 0x82, 0x9a,
	0x81, 0xba, 0x14, 0x17, 0x76, 0xa1, 0x14, 0x6f, 0x65, 0x04, 0x17, 0x6e, 0xca, 0x24, 0x1e, 0xd3,
	0x60, 0x9b, 0x09, 0x33, 0x53, 0xb1, 0x3e, 0x85, 0x4f, 0xe4, 0xda, 0x65, 0x97, 0x2e, 0xa5, 0x7d,
	0x11, 0xc9, 0x4c, 0x2e, 0xd8, 0x42, 0x70, 0x77, 0x98, 0xf3, 0x7f, 0x5f, 0x72, 0xce, 0x0c, 0xda,
	0x61, 0xae, 0x00, 0xfe, 0x0a, 0x1c, 0xfb, 0x10, 0x82, 0x08, 0x84, 0x13, 0x71, 0x26, 0x99, 0xd9,
	0x78, 0x07, 0x49, 0xbd, 0x01, 0x0d, 0x42, 0x47, 0x55, 0x8c, 0x83, 0x93, 0x46, 0xeb, 0x5b, 0x3e,
	0xf3, 0x99, 0xca, 0xe1, 0xb8, 0xd2, 0x48, 0x7d, 0x37, 0x53, 0xa5, 0x45, 0xd2, 0xd8, 0xce, 0x1a,
	0x2e, 0x1d, 0x0e, 0x99, 0x5c, 0x3a, 0x8e, 0x28, 0xa7, 0xa3, 0xe4, 0xcb, 0xf5, 0x46, 0x76, 0x1c,
	0xb2, 0x27, 0xe8, 0x53, 0xcf, 0x63, 0xe3, 0x30, 0x65, 0xf6, 0x72, 0x06, 0xf8, 0x28, 0x10, 0x22,
	0x60, 0x61, 0xff, 0x79, 0x48, 0x7d, 0xb1, 0x24, 0x7d, 0x81, 0x89, 0x0f, 0xa1, 0x3e, 0xde, 0xff,
	0x2c, 0xa3, 0xcd, 0x4b, 0x3d, 0xe0, 0xbd, 0xa4, 0x12, 0xcc, 0x33, 0xb4, 0xae, 0x7f, 0x46, 0x58,
	0x46, 0xb3, 0xdc, 0xaa, 0xb6, 0x0f, 0x9c, 0x82, 0x89, 0x9d, 0x8e, 0xca, 0x92, 0x94, 0x31, 0xbb,
	0x68, 0x23, 0xed, 0x09, 0x6b, 0x45, 0x09, 0x0e, 0x0b, 0x05, 0x77, 0x49, 0x71, 0x43, 0xa3, 0x08,
	0x38, 0xc9, 0x69, 0x93, 0xa0, 0x5a, 0x3c, 0xe8, 0xb9, 0x9e, 0xf3, 0x3a, 0x10, 0xd2, 0x2a, 0x2b,
	0x61, 0xab, 0x50, 0x78, 0x9b, 0x33, 0x64, 0x51, 0x60, 0x3e, 0xa0, 0x5a, 0xbe, 0x9f, 0x8b, 0x78,
	0x3d, 0xd6, 0x6a, 0xd3, 0x68, 0x55, 0xdb, 0x47, 0x85, 0xce, 0xde, 0x5f, 0x86, 0x2c, 0x4a, 0xcc,
	0x53, 0x54, 0xd1, 0x77, 0x65, 0xad, 0x29, 0x5d, 0xf1, 0xd2, 0x7a, 0x2a, 0x4a, 0x12, 0x24, 0x86,
	0xf5, 0x9d, 0x58, 0x95, 0x7f, 0xc0, 0x57, 0x2a, 0x4a, 0x12, 0xa4, 0xd3, 0xfd, 0x9a, 0xd9, 0xc6,
	0x74, 0x66, 0x1b, 0x3f, 0x33, 0xdb, 0xf8, 0x98, 0xdb, 0xa5, 0xe9, 0xdc, 0x2e, 0x7d, 0xcf, 0xed,
	0xd2, 0x23, 0xf6, 0x03, 0x39, 0x18, 0xbb, 0x8e, 0xc7, 0x46, 0x38, 0xd6, 0x1c, 0x2b, 0x23, 0x4e,
	0x8d, 0xf8, 0x2d, 0x7b, 0x8e, 0x58, 0x4e, 0x22, 0x10, 0x6e, 0x45, 0x3d, 0x89, 0x93, 0xdf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x0a, 0x4e, 0xe0, 0xbc, 0xfb, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Keygen != nil {
		{
			size, err := m.Keygen.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.Params != nil {
		{
			size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.PermissionFlags != nil {
		{
			size, err := m.PermissionFlags.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.NodeAccountList) > 0 {
		for iNdEx := len(m.NodeAccountList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.NodeAccountList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Observers) > 0 {
		for iNdEx := len(m.Observers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Observers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Ballots) > 0 {
		for iNdEx := len(m.Ballots) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Ballots[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ballots) > 0 {
		for _, e := range m.Ballots {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Observers) > 0 {
		for _, e := range m.Observers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.NodeAccountList) > 0 {
		for _, e := range m.NodeAccountList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if m.PermissionFlags != nil {
		l = m.PermissionFlags.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Params != nil {
		l = m.Params.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Keygen != nil {
		l = m.Keygen.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ballots", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ballots = append(m.Ballots, &Ballot{})
			if err := m.Ballots[len(m.Ballots)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Observers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Observers = append(m.Observers, &ObserverMapper{})
			if err := m.Observers[len(m.Observers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NodeAccountList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NodeAccountList = append(m.NodeAccountList, &NodeAccount{})
			if err := m.NodeAccountList[len(m.NodeAccountList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PermissionFlags", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PermissionFlags == nil {
				m.PermissionFlags = &PermissionFlags{}
			}
			if err := m.PermissionFlags.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Params == nil {
				m.Params = &Params{}
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keygen", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Keygen == nil {
				m.Keygen = &Keygen{}
			}
			if err := m.Keygen.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
