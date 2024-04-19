// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lightclient/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	proofs "github.com/zeta-chain/zetacore/pkg/proofs"
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

// GenesisState defines the lightclient module's genesis state.
type GenesisState struct {
	BlockHeaders      []proofs.BlockHeader `protobuf:"bytes,1,rep,name=block_headers,json=blockHeaders,proto3" json:"block_headers"`
	ChainStates       []ChainState         `protobuf:"bytes,2,rep,name=chain_states,json=chainStates,proto3" json:"chain_states"`
	VerificationFlags VerificationFlags    `protobuf:"bytes,3,opt,name=verification_flags,json=verificationFlags,proto3" json:"verification_flags"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_645b5300b371cd43, []int{0}
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

func (m *GenesisState) GetBlockHeaders() []proofs.BlockHeader {
	if m != nil {
		return m.BlockHeaders
	}
	return nil
}

func (m *GenesisState) GetChainStates() []ChainState {
	if m != nil {
		return m.ChainStates
	}
	return nil
}

func (m *GenesisState) GetVerificationFlags() VerificationFlags {
	if m != nil {
		return m.VerificationFlags
	}
	return VerificationFlags{}
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "lightclient.GenesisState")
}

func init() { proto.RegisterFile("lightclient/genesis.proto", fileDescriptor_645b5300b371cd43) }

var fileDescriptor_645b5300b371cd43 = []byte{
	// 307 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0xc1, 0x4a, 0xf4, 0x30,
	0x14, 0x85, 0x9b, 0x7f, 0x7e, 0x5c, 0xb4, 0xe3, 0xc2, 0x22, 0x4c, 0x2d, 0x18, 0x07, 0x71, 0x31,
	0x1b, 0x1b, 0xa8, 0x2f, 0x20, 0x15, 0xd4, 0x85, 0x2b, 0x07, 0x5c, 0xb8, 0x29, 0x69, 0x4c, 0xd3,
	0xd0, 0xda, 0x94, 0x26, 0x0e, 0xea, 0x53, 0xf8, 0x58, 0xb3, 0x9c, 0xa5, 0x2b, 0x91, 0x76, 0xe5,
	0x5b, 0x48, 0xd3, 0xa0, 0x91, 0x59, 0xe5, 0x72, 0xbf, 0x73, 0xcf, 0x3d, 0xb9, 0xee, 0x41, 0xc5,
	0x59, 0xa1, 0x48, 0xc5, 0x69, 0xad, 0x10, 0xa3, 0x35, 0x95, 0x5c, 0x46, 0x4d, 0x2b, 0x94, 0xf0,
	0x3d, 0x0b, 0x85, 0xfb, 0x4c, 0x30, 0xa1, 0xfb, 0x68, 0xa8, 0x46, 0x49, 0x78, 0x68, 0x4f, 0x93,
	0x02, 0xf3, 0x3a, 0x95, 0x0a, 0x2b, 0x6a, 0xf0, 0x89, 0x8d, 0x57, 0xb4, 0xe5, 0x39, 0x27, 0x58,
	0x71, 0x51, 0xa7, 0x79, 0x85, 0x99, 0xd9, 0x13, 0xce, 0x9a, 0x92, 0xa1, 0xa6, 0x15, 0x22, 0x97,
	0xe6, 0x19, 0xc1, 0xf1, 0x17, 0x70, 0xa7, 0x57, 0x63, 0xa4, 0xe5, 0xe0, 0xea, 0x27, 0xee, 0x6e,
	0x56, 0x09, 0x52, 0xa6, 0x05, 0xc5, 0x0f, 0xb4, 0x95, 0x01, 0x98, 0x4f, 0x16, 0x5e, 0x3c, 0x8b,
	0x9a, 0x92, 0x45, 0x66, 0x34, 0x19, 0x04, 0xd7, 0x9a, 0x27, 0xff, 0xd7, 0x1f, 0x47, 0xce, 0xed,
	0x34, 0xfb, 0x6d, 0x49, 0xff, 0xdc, 0x9d, 0x5a, 0x41, 0x65, 0xf0, 0xcf, 0x58, 0x58, 0x51, 0xa3,
	0x8b, 0x41, 0xa0, 0x57, 0x1a, 0x0b, 0x8f, 0xfc, 0x74, 0xa4, 0xbf, 0x74, 0xfd, 0xed, 0xbf, 0x04,
	0x93, 0x39, 0x58, 0x78, 0x31, 0xfc, 0xe3, 0x73, 0x67, 0xc9, 0x2e, 0x07, 0x95, 0xb1, 0xdb, 0x5b,
	0x6d, 0x81, 0x9b, 0x75, 0x07, 0xc1, 0xa6, 0x83, 0xe0, 0xb3, 0x83, 0xe0, 0xad, 0x87, 0xce, 0xa6,
	0x87, 0xce, 0x7b, 0x0f, 0x9d, 0xfb, 0x98, 0x71, 0x55, 0x3c, 0x65, 0x11, 0x11, 0x8f, 0xe8, 0x95,
	0x2a, 0x7c, 0xaa, 0xb3, 0xe8, 0x92, 0x88, 0x96, 0xa2, 0x67, 0x64, 0x5f, 0x59, 0xbd, 0x34, 0x54,
	0x66, 0x3b, 0xfa, 0x80, 0x67, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x9e, 0xe5, 0xcf, 0x01, 0xde,
	0x01, 0x00, 0x00,
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
	{
		size, err := m.VerificationFlags.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ChainStates) > 0 {
		for iNdEx := len(m.ChainStates) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ChainStates[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.BlockHeaders) > 0 {
		for iNdEx := len(m.BlockHeaders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.BlockHeaders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.BlockHeaders) > 0 {
		for _, e := range m.BlockHeaders {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ChainStates) > 0 {
		for _, e := range m.ChainStates {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.VerificationFlags.Size()
	n += 1 + l + sovGenesis(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeaders", wireType)
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
			m.BlockHeaders = append(m.BlockHeaders, proofs.BlockHeader{})
			if err := m.BlockHeaders[len(m.BlockHeaders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainStates", wireType)
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
			m.ChainStates = append(m.ChainStates, ChainState{})
			if err := m.ChainStates[len(m.ChainStates)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VerificationFlags", wireType)
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
			if err := m.VerificationFlags.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
