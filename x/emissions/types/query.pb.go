// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: emissions/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e578782beb6ef82, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params holds all the parameters of this module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e578782beb6ef82, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

type QueryListBalancesRequest struct {
}

func (m *QueryListBalancesRequest) Reset()         { *m = QueryListBalancesRequest{} }
func (m *QueryListBalancesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryListBalancesRequest) ProtoMessage()    {}
func (*QueryListBalancesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e578782beb6ef82, []int{2}
}
func (m *QueryListBalancesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryListBalancesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryListBalancesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryListBalancesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryListBalancesRequest.Merge(m, src)
}
func (m *QueryListBalancesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryListBalancesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryListBalancesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryListBalancesRequest proto.InternalMessageInfo

type QueryListBalancesResponse struct {
	Trackers              []*EmissionTracker `protobuf:"bytes,1,rep,name=trackers,proto3" json:"trackers,omitempty"`
	EmissionModuleAddress string             `protobuf:"bytes,2,opt,name=emission_module_address,json=emissionModuleAddress,proto3" json:"emission_module_address,omitempty"`
}

func (m *QueryListBalancesResponse) Reset()         { *m = QueryListBalancesResponse{} }
func (m *QueryListBalancesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryListBalancesResponse) ProtoMessage()    {}
func (*QueryListBalancesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6e578782beb6ef82, []int{3}
}
func (m *QueryListBalancesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryListBalancesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryListBalancesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryListBalancesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryListBalancesResponse.Merge(m, src)
}
func (m *QueryListBalancesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryListBalancesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryListBalancesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryListBalancesResponse proto.InternalMessageInfo

func (m *QueryListBalancesResponse) GetTrackers() []*EmissionTracker {
	if m != nil {
		return m.Trackers
	}
	return nil
}

func (m *QueryListBalancesResponse) GetEmissionModuleAddress() string {
	if m != nil {
		return m.EmissionModuleAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "zetachain.zetacore.emissions.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "zetachain.zetacore.emissions.QueryParamsResponse")
	proto.RegisterType((*QueryListBalancesRequest)(nil), "zetachain.zetacore.emissions.QueryListBalancesRequest")
	proto.RegisterType((*QueryListBalancesResponse)(nil), "zetachain.zetacore.emissions.QueryListBalancesResponse")
}

func init() { proto.RegisterFile("emissions/query.proto", fileDescriptor_6e578782beb6ef82) }

var fileDescriptor_6e578782beb6ef82 = []byte{
	// 442 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0x41, 0xcb, 0xd3, 0x30,
	0x1c, 0xc6, 0x9b, 0xa9, 0x43, 0xf3, 0x7a, 0x8a, 0xaf, 0x3a, 0xcb, 0x4b, 0x1d, 0x45, 0x71, 0x88,
	0x6b, 0xde, 0x4d, 0x99, 0x67, 0x0b, 0x1e, 0x14, 0x05, 0x2d, 0x5e, 0xf4, 0x32, 0xd2, 0x2e, 0x74,
	0xc1, 0xb6, 0xe9, 0x9a, 0x54, 0x9c, 0x47, 0x3f, 0x81, 0xe0, 0x55, 0xfc, 0x16, 0x7e, 0x87, 0x1d,
	0x07, 0x5e, 0x3c, 0x0d, 0xd9, 0xfc, 0x20, 0xd2, 0x24, 0xeb, 0xa6, 0x8e, 0xf9, 0xee, 0xf6, 0xa7,
	0xff, 0xe7, 0xf7, 0x3c, 0x4f, 0x92, 0xc2, 0xab, 0x34, 0x65, 0x42, 0x30, 0x9e, 0x09, 0x3c, 0x29,
	0x69, 0x31, 0xf5, 0xf2, 0x82, 0x4b, 0x8e, 0x4e, 0x3e, 0x50, 0x49, 0xa2, 0x31, 0x61, 0x99, 0xa7,
	0x26, 0x5e, 0x50, 0xaf, 0x56, 0xda, 0xc7, 0x31, 0x8f, 0xb9, 0x12, 0xe2, 0x6a, 0xd2, 0x8c, 0x7d,
	0x12, 0x73, 0x1e, 0x27, 0x14, 0x93, 0x9c, 0x61, 0x92, 0x65, 0x5c, 0x12, 0x59, 0xa9, 0xcd, 0xf6,
	0x6e, 0xc4, 0x45, 0xca, 0x05, 0x0e, 0x89, 0xa0, 0x3a, 0x0a, 0xbf, 0xeb, 0x85, 0x54, 0x92, 0x1e,
	0xce, 0x49, 0xcc, 0x32, 0x25, 0x36, 0xda, 0x6b, 0x9b, 0x52, 0x39, 0x29, 0x48, 0xba, 0xf6, 0x68,
	0x6f, 0xbe, 0xaf, 0xa7, 0xa1, 0x2c, 0x48, 0xf4, 0x96, 0x16, 0x5a, 0xe1, 0x1e, 0x43, 0xf4, 0xb2,
	0xf2, 0x7e, 0xa1, 0xb0, 0x80, 0x4e, 0x4a, 0x2a, 0xa4, 0xfb, 0x1a, 0x5e, 0xf9, 0xe3, 0xab, 0xc8,
	0x79, 0x26, 0x28, 0xf2, 0x61, 0x53, 0xdb, 0xb7, 0x40, 0x1b, 0x74, 0x8e, 0xfa, 0xb7, 0xbc, 0x7d,
	0xa7, 0xf6, 0x34, 0xed, 0x9f, 0x9f, 0x2d, 0x6e, 0x5a, 0x81, 0x21, 0x5d, 0x1b, 0xb6, 0x94, 0xf5,
	0x33, 0x26, 0xa4, 0x4f, 0x12, 0x92, 0x45, 0xb4, 0x8e, 0xfd, 0x0a, 0xe0, 0x8d, 0x1d, 0x4b, 0x93,
	0xfe, 0x04, 0x5e, 0x34, 0xdd, 0xab, 0xfc, 0x73, 0x9d, 0xa3, 0x7e, 0x77, 0x7f, 0xfe, 0x63, 0x33,
	0xbd, 0xd2, 0x54, 0x50, 0xe3, 0x68, 0x00, 0xaf, 0xd7, 0xf7, 0x91, 0xf2, 0x51, 0x99, 0xd0, 0x21,
	0x19, 0x8d, 0x0a, 0x2a, 0x44, 0xab, 0xd1, 0x06, 0x9d, 0x4b, 0x41, 0xfd, 0xca, 0xcf, 0xd5, 0xf6,
	0x91, 0x5e, 0xf6, 0x17, 0x0d, 0x78, 0x41, 0x15, 0x44, 0x5f, 0x00, 0x6c, 0xea, 0xf3, 0xa1, 0xd3,
	0xfd, 0x2d, 0xfe, 0xbd, 0x5e, 0xbb, 0x77, 0x00, 0xa1, 0x0f, 0xef, 0x76, 0x3f, 0x7e, 0xff, 0xf5,
	0xb9, 0x71, 0x07, 0xdd, 0xc6, 0x15, 0xd0, 0x55, 0x2c, 0x5e, 0xb3, 0xf8, 0xef, 0xe7, 0x47, 0xdf,
	0x00, 0xbc, 0xbc, 0x7d, 0x89, 0x68, 0x70, 0x86, 0xc8, 0x1d, 0x4f, 0x62, 0x3f, 0x3c, 0x98, 0x33,
	0x85, 0x1f, 0xa8, 0xc2, 0x1e, 0xba, 0xf7, 0x9f, 0xc2, 0x09, 0x13, 0x72, 0x18, 0x1a, 0xda, 0x7f,
	0x3a, 0x5b, 0x3a, 0x60, 0xbe, 0x74, 0xc0, 0xcf, 0xa5, 0x03, 0x3e, 0xad, 0x1c, 0x6b, 0xbe, 0x72,
	0xac, 0x1f, 0x2b, 0xc7, 0x7a, 0x73, 0x1a, 0x33, 0x39, 0x2e, 0x43, 0x2f, 0xe2, 0xe9, 0x4e, 0xc7,
	0xf7, 0x5b, 0x9e, 0x72, 0x9a, 0x53, 0x11, 0x36, 0xd5, 0x1f, 0x7e, 0xff, 0x77, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x07, 0x77, 0xaf, 0x30, 0xb2, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of ListBalances items.
	ListBalances(ctx context.Context, in *QueryListBalancesRequest, opts ...grpc.CallOption) (*QueryListBalancesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/zetachain.zetacore.emissions.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ListBalances(ctx context.Context, in *QueryListBalancesRequest, opts ...grpc.CallOption) (*QueryListBalancesResponse, error) {
	out := new(QueryListBalancesResponse)
	err := c.cc.Invoke(ctx, "/zetachain.zetacore.emissions.Query/ListBalances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of ListBalances items.
	ListBalances(context.Context, *QueryListBalancesRequest) (*QueryListBalancesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) ListBalances(ctx context.Context, req *QueryListBalancesRequest) (*QueryListBalancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListBalances not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zetachain.zetacore.emissions.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ListBalances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryListBalancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ListBalances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zetachain.zetacore.emissions.Query/ListBalances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ListBalances(ctx, req.(*QueryListBalancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "zetachain.zetacore.emissions.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "ListBalances",
			Handler:    _Query_ListBalances_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "emissions/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryListBalancesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryListBalancesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryListBalancesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryListBalancesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryListBalancesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryListBalancesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.EmissionModuleAddress) > 0 {
		i -= len(m.EmissionModuleAddress)
		copy(dAtA[i:], m.EmissionModuleAddress)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.EmissionModuleAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Trackers) > 0 {
		for iNdEx := len(m.Trackers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Trackers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryListBalancesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryListBalancesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Trackers) > 0 {
		for _, e := range m.Trackers {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	l = len(m.EmissionModuleAddress)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryListBalancesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryListBalancesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryListBalancesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryListBalancesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryListBalancesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryListBalancesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Trackers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Trackers = append(m.Trackers, &EmissionTracker{})
			if err := m.Trackers[len(m.Trackers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EmissionModuleAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EmissionModuleAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
