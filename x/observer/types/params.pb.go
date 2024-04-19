// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: observer/params.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	chains "github.com/zeta-chain/zetacore/pkg/chains"
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

type ChainParamsList struct {
	ChainParams []*ChainParams `protobuf:"bytes,1,rep,name=chain_params,json=chainParams,proto3" json:"chain_params,omitempty"`
}

func (m *ChainParamsList) Reset()         { *m = ChainParamsList{} }
func (m *ChainParamsList) String() string { return proto.CompactTextString(m) }
func (*ChainParamsList) ProtoMessage()    {}
func (*ChainParamsList) Descriptor() ([]byte, []int) {
	return fileDescriptor_4542fa62877488a1, []int{0}
}
func (m *ChainParamsList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainParamsList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainParamsList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainParamsList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainParamsList.Merge(m, src)
}
func (m *ChainParamsList) XXX_Size() int {
	return m.Size()
}
func (m *ChainParamsList) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainParamsList.DiscardUnknown(m)
}

var xxx_messageInfo_ChainParamsList proto.InternalMessageInfo

func (m *ChainParamsList) GetChainParams() []*ChainParams {
	if m != nil {
		return m.ChainParams
	}
	return nil
}

type ChainParams struct {
	ChainId                     int64                                  `protobuf:"varint,11,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	ConfirmationCount           uint64                                 `protobuf:"varint,1,opt,name=confirmation_count,json=confirmationCount,proto3" json:"confirmation_count,omitempty"`
	GasPriceTicker              uint64                                 `protobuf:"varint,2,opt,name=gas_price_ticker,json=gasPriceTicker,proto3" json:"gas_price_ticker,omitempty"`
	InTxTicker                  uint64                                 `protobuf:"varint,3,opt,name=in_tx_ticker,json=inTxTicker,proto3" json:"in_tx_ticker,omitempty"`
	OutTxTicker                 uint64                                 `protobuf:"varint,4,opt,name=out_tx_ticker,json=outTxTicker,proto3" json:"out_tx_ticker,omitempty"`
	WatchUtxoTicker             uint64                                 `protobuf:"varint,5,opt,name=watch_utxo_ticker,json=watchUtxoTicker,proto3" json:"watch_utxo_ticker,omitempty"`
	ZetaTokenContractAddress    string                                 `protobuf:"bytes,8,opt,name=zeta_token_contract_address,json=zetaTokenContractAddress,proto3" json:"zeta_token_contract_address,omitempty"`
	ConnectorContractAddress    string                                 `protobuf:"bytes,9,opt,name=connector_contract_address,json=connectorContractAddress,proto3" json:"connector_contract_address,omitempty"`
	Erc20CustodyContractAddress string                                 `protobuf:"bytes,10,opt,name=erc20_custody_contract_address,json=erc20CustodyContractAddress,proto3" json:"erc20_custody_contract_address,omitempty"`
	OutboundTxScheduleInterval  int64                                  `protobuf:"varint,12,opt,name=outbound_tx_schedule_interval,json=outboundTxScheduleInterval,proto3" json:"outbound_tx_schedule_interval,omitempty"`
	OutboundTxScheduleLookahead int64                                  `protobuf:"varint,13,opt,name=outbound_tx_schedule_lookahead,json=outboundTxScheduleLookahead,proto3" json:"outbound_tx_schedule_lookahead,omitempty"`
	BallotThreshold             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=ballot_threshold,json=ballotThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"ballot_threshold"`
	MinObserverDelegation       github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,15,opt,name=min_observer_delegation,json=minObserverDelegation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_observer_delegation"`
	IsSupported                 bool                                   `protobuf:"varint,16,opt,name=is_supported,json=isSupported,proto3" json:"is_supported,omitempty"`
}

func (m *ChainParams) Reset()         { *m = ChainParams{} }
func (m *ChainParams) String() string { return proto.CompactTextString(m) }
func (*ChainParams) ProtoMessage()    {}
func (*ChainParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_4542fa62877488a1, []int{1}
}
func (m *ChainParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ChainParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ChainParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ChainParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChainParams.Merge(m, src)
}
func (m *ChainParams) XXX_Size() int {
	return m.Size()
}
func (m *ChainParams) XXX_DiscardUnknown() {
	xxx_messageInfo_ChainParams.DiscardUnknown(m)
}

var xxx_messageInfo_ChainParams proto.InternalMessageInfo

func (m *ChainParams) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *ChainParams) GetConfirmationCount() uint64 {
	if m != nil {
		return m.ConfirmationCount
	}
	return 0
}

func (m *ChainParams) GetGasPriceTicker() uint64 {
	if m != nil {
		return m.GasPriceTicker
	}
	return 0
}

func (m *ChainParams) GetInTxTicker() uint64 {
	if m != nil {
		return m.InTxTicker
	}
	return 0
}

func (m *ChainParams) GetOutTxTicker() uint64 {
	if m != nil {
		return m.OutTxTicker
	}
	return 0
}

func (m *ChainParams) GetWatchUtxoTicker() uint64 {
	if m != nil {
		return m.WatchUtxoTicker
	}
	return 0
}

func (m *ChainParams) GetZetaTokenContractAddress() string {
	if m != nil {
		return m.ZetaTokenContractAddress
	}
	return ""
}

func (m *ChainParams) GetConnectorContractAddress() string {
	if m != nil {
		return m.ConnectorContractAddress
	}
	return ""
}

func (m *ChainParams) GetErc20CustodyContractAddress() string {
	if m != nil {
		return m.Erc20CustodyContractAddress
	}
	return ""
}

func (m *ChainParams) GetOutboundTxScheduleInterval() int64 {
	if m != nil {
		return m.OutboundTxScheduleInterval
	}
	return 0
}

func (m *ChainParams) GetOutboundTxScheduleLookahead() int64 {
	if m != nil {
		return m.OutboundTxScheduleLookahead
	}
	return 0
}

func (m *ChainParams) GetIsSupported() bool {
	if m != nil {
		return m.IsSupported
	}
	return false
}

// Deprecated(v13): Use ChainParamsList
type ObserverParams struct {
	Chain                 *chains.Chain                          `protobuf:"bytes,1,opt,name=chain,proto3" json:"chain,omitempty"`
	BallotThreshold       github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=ballot_threshold,json=ballotThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"ballot_threshold"`
	MinObserverDelegation github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=min_observer_delegation,json=minObserverDelegation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_observer_delegation"`
	IsSupported           bool                                   `protobuf:"varint,5,opt,name=is_supported,json=isSupported,proto3" json:"is_supported,omitempty"`
}

func (m *ObserverParams) Reset()         { *m = ObserverParams{} }
func (m *ObserverParams) String() string { return proto.CompactTextString(m) }
func (*ObserverParams) ProtoMessage()    {}
func (*ObserverParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_4542fa62877488a1, []int{2}
}
func (m *ObserverParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ObserverParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ObserverParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ObserverParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ObserverParams.Merge(m, src)
}
func (m *ObserverParams) XXX_Size() int {
	return m.Size()
}
func (m *ObserverParams) XXX_DiscardUnknown() {
	xxx_messageInfo_ObserverParams.DiscardUnknown(m)
}

var xxx_messageInfo_ObserverParams proto.InternalMessageInfo

func (m *ObserverParams) GetChain() *chains.Chain {
	if m != nil {
		return m.Chain
	}
	return nil
}

func (m *ObserverParams) GetIsSupported() bool {
	if m != nil {
		return m.IsSupported
	}
	return false
}

func init() {
	proto.RegisterType((*ChainParamsList)(nil), "observer.ChainParamsList")
	proto.RegisterType((*ChainParams)(nil), "observer.ChainParams")
	proto.RegisterType((*ObserverParams)(nil), "observer.ObserverParams")
}

func init() { proto.RegisterFile("observer/params.proto", fileDescriptor_4542fa62877488a1) }

var fileDescriptor_4542fa62877488a1 = []byte{
	// 638 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xc1, 0x4e, 0xdb, 0x4c,
	0x10, 0xc7, 0x63, 0x02, 0x7c, 0x61, 0x1d, 0x08, 0x58, 0x1f, 0xc2, 0x0d, 0xaa, 0x49, 0x73, 0x68,
	0xa3, 0x4a, 0xd8, 0x15, 0xbd, 0xf4, 0xd0, 0x1e, 0x20, 0x5c, 0x50, 0x91, 0x8a, 0x4c, 0x7a, 0x68,
	0x2f, 0xab, 0xcd, 0x7a, 0xb1, 0x57, 0x71, 0x3c, 0xd6, 0xee, 0x9a, 0x86, 0x3e, 0x45, 0xdf, 0xa0,
	0x87, 0xbe, 0x0c, 0x47, 0x8e, 0x55, 0x0f, 0xa8, 0x22, 0x2f, 0x52, 0x79, 0x6d, 0xa7, 0x11, 0xe1,
	0x54, 0x89, 0x93, 0x77, 0xe7, 0xff, 0x9b, 0xff, 0x6c, 0x26, 0x3b, 0x8b, 0xb6, 0x61, 0x28, 0x99,
	0xb8, 0x64, 0xc2, 0x4b, 0x89, 0x20, 0x63, 0xe9, 0xa6, 0x02, 0x14, 0x58, 0x8d, 0x2a, 0xdc, 0xfe,
	0x3f, 0x84, 0x10, 0x74, 0xd0, 0xcb, 0x57, 0x85, 0xde, 0xde, 0x99, 0xa5, 0x55, 0x8b, 0x4a, 0x48,
	0x47, 0xa1, 0x47, 0x23, 0xc2, 0x13, 0x59, 0x7e, 0x0a, 0xa1, 0xfb, 0x1e, 0xb5, 0xfa, 0xf9, 0xfe,
	0x4c, 0x97, 0x39, 0xe5, 0x52, 0x59, 0x6f, 0x50, 0x53, 0x23, 0xb8, 0x28, 0x6d, 0x1b, 0x9d, 0x7a,
	0xcf, 0x3c, 0xd8, 0x76, 0x67, 0x96, 0x73, 0x09, 0xbe, 0x49, 0xff, 0x6e, 0xba, 0x3f, 0x56, 0x91,
	0x39, 0x27, 0x5a, 0x4f, 0x50, 0xa3, 0x70, 0xe2, 0x81, 0x6d, 0x76, 0x8c, 0x5e, 0xdd, 0xff, 0x4f,
	0xef, 0x4f, 0x02, 0x6b, 0x1f, 0x59, 0x14, 0x92, 0x0b, 0x2e, 0xc6, 0x44, 0x71, 0x48, 0x30, 0x85,
	0x2c, 0x51, 0xb6, 0xd1, 0x31, 0x7a, 0xcb, 0xfe, 0xd6, 0xbc, 0xd2, 0xcf, 0x05, 0xab, 0x87, 0x36,
	0x43, 0x22, 0x71, 0x2a, 0x38, 0x65, 0x58, 0x71, 0x3a, 0x62, 0xc2, 0x5e, 0xd2, 0xf0, 0x46, 0x48,
	0xe4, 0x59, 0x1e, 0x1e, 0xe8, 0xa8, 0xd5, 0x41, 0x4d, 0x9e, 0x60, 0x35, 0xa9, 0xa8, 0xba, 0xa6,
	0x10, 0x4f, 0x06, 0x93, 0x92, 0xe8, 0xa2, 0x75, 0xc8, 0xd4, 0x1c, 0xb2, 0xac, 0x11, 0x13, 0x32,
	0x35, 0x63, 0x5e, 0xa2, 0xad, 0x2f, 0x44, 0xd1, 0x08, 0x67, 0x6a, 0x02, 0x15, 0xb7, 0xa2, 0xb9,
	0x96, 0x16, 0x3e, 0xaa, 0x09, 0x94, 0xec, 0x3b, 0xb4, 0xfb, 0x95, 0x29, 0x82, 0x15, 0x8c, 0x58,
	0xfe, 0x43, 0x12, 0x25, 0x08, 0x55, 0x98, 0x04, 0x81, 0x60, 0x52, 0xda, 0x8d, 0x8e, 0xd1, 0x5b,
	0xf3, 0xed, 0x1c, 0x19, 0xe4, 0x44, 0xbf, 0x04, 0x0e, 0x0b, 0xdd, 0x7a, 0x8b, 0xda, 0x14, 0x92,
	0x84, 0x51, 0x05, 0x62, 0x31, 0x7b, 0xad, 0xc8, 0x9e, 0x11, 0xf7, 0xb3, 0xfb, 0xc8, 0x61, 0x82,
	0x1e, 0xbc, 0xc2, 0x34, 0x93, 0x0a, 0x82, 0xab, 0x45, 0x07, 0xa4, 0x1d, 0x76, 0x35, 0xd5, 0x2f,
	0xa0, 0xfb, 0x26, 0x87, 0xe8, 0x29, 0x64, 0x6a, 0x08, 0x59, 0x12, 0xe4, 0x6d, 0x91, 0x34, 0x62,
	0x41, 0x16, 0x33, 0xcc, 0x13, 0xc5, 0xc4, 0x25, 0x89, 0xed, 0xa6, 0xfe, 0xf3, 0xda, 0x15, 0x34,
	0x98, 0x9c, 0x97, 0xc8, 0x49, 0x49, 0xe4, 0xe7, 0x78, 0xd0, 0x22, 0x06, 0x18, 0x91, 0x88, 0x91,
	0xc0, 0x5e, 0xd7, 0x1e, 0xbb, 0x8b, 0x1e, 0xa7, 0x15, 0x62, 0x7d, 0x42, 0x9b, 0x43, 0x12, 0xc7,
	0xa0, 0xb0, 0x8a, 0x04, 0x93, 0x11, 0xc4, 0x81, 0xbd, 0x91, 0x1f, 0xff, 0xc8, 0xbd, 0xbe, 0xdd,
	0xab, 0xfd, 0xba, 0xdd, 0x7b, 0x1e, 0x72, 0x15, 0x65, 0x43, 0x97, 0xc2, 0xd8, 0xa3, 0x20, 0xc7,
	0x20, 0xcb, 0xcf, 0xbe, 0x0c, 0x46, 0x9e, 0xba, 0x4a, 0x99, 0x74, 0x8f, 0x19, 0xf5, 0x5b, 0x85,
	0xcf, 0xa0, 0xb2, 0xb1, 0x2e, 0xd0, 0xce, 0x98, 0x27, 0xb8, 0xba, 0xc3, 0x38, 0x60, 0x31, 0x0b,
	0xf5, 0x05, 0xb3, 0x5b, 0xff, 0x54, 0x61, 0x7b, 0xcc, 0x93, 0x0f, 0xa5, 0xdb, 0xf1, 0xcc, 0xcc,
	0x7a, 0x86, 0x9a, 0x5c, 0x62, 0x99, 0xa5, 0x29, 0x08, 0xc5, 0x02, 0x7b, 0xb3, 0x63, 0xf4, 0x1a,
	0xbe, 0xc9, 0xe5, 0x79, 0x15, 0xea, 0x7e, 0x5f, 0x42, 0x1b, 0x55, 0x66, 0x39, 0x28, 0x2f, 0xd0,
	0x8a, 0x1e, 0x0c, 0x3d, 0x00, 0xe6, 0xc1, 0x96, 0x9b, 0x8e, 0x42, 0xb7, 0x9c, 0x53, 0x3d, 0x50,
	0x7e, 0xa1, 0x3f, 0xd8, 0xa1, 0xfa, 0xa3, 0x77, 0x68, 0xf9, 0x31, 0x3b, 0xb4, 0xb2, 0xd0, 0xa1,
	0xa3, 0x93, 0xeb, 0x3b, 0xc7, 0xb8, 0xb9, 0x73, 0x8c, 0xdf, 0x77, 0x8e, 0xf1, 0x6d, 0xea, 0xd4,
	0x6e, 0xa6, 0x4e, 0xed, 0xe7, 0xd4, 0xa9, 0x7d, 0xf6, 0xe6, 0x6a, 0xe7, 0x13, 0xb5, 0xaf, 0xdb,
	0xa2, 0x97, 0x14, 0x04, 0xf3, 0x26, 0xb3, 0x87, 0xaf, 0x38, 0xc8, 0x70, 0x55, 0x3f, 0x73, 0xaf,
	0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd9, 0x5c, 0x62, 0xf2, 0x51, 0x05, 0x00, 0x00,
}

func (m *ChainParamsList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainParamsList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainParamsList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ChainParams) > 0 {
		for iNdEx := len(m.ChainParams) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ChainParams[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *ChainParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ChainParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ChainParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsSupported {
		i--
		if m.IsSupported {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x80
	}
	{
		size := m.MinObserverDelegation.Size()
		i -= size
		if _, err := m.MinObserverDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	{
		size := m.BallotThreshold.Size()
		i -= size
		if _, err := m.BallotThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x72
	if m.OutboundTxScheduleLookahead != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OutboundTxScheduleLookahead))
		i--
		dAtA[i] = 0x68
	}
	if m.OutboundTxScheduleInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OutboundTxScheduleInterval))
		i--
		dAtA[i] = 0x60
	}
	if m.ChainId != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x58
	}
	if len(m.Erc20CustodyContractAddress) > 0 {
		i -= len(m.Erc20CustodyContractAddress)
		copy(dAtA[i:], m.Erc20CustodyContractAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Erc20CustodyContractAddress)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.ConnectorContractAddress) > 0 {
		i -= len(m.ConnectorContractAddress)
		copy(dAtA[i:], m.ConnectorContractAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ConnectorContractAddress)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.ZetaTokenContractAddress) > 0 {
		i -= len(m.ZetaTokenContractAddress)
		copy(dAtA[i:], m.ZetaTokenContractAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ZetaTokenContractAddress)))
		i--
		dAtA[i] = 0x42
	}
	if m.WatchUtxoTicker != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.WatchUtxoTicker))
		i--
		dAtA[i] = 0x28
	}
	if m.OutTxTicker != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OutTxTicker))
		i--
		dAtA[i] = 0x20
	}
	if m.InTxTicker != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.InTxTicker))
		i--
		dAtA[i] = 0x18
	}
	if m.GasPriceTicker != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.GasPriceTicker))
		i--
		dAtA[i] = 0x10
	}
	if m.ConfirmationCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ConfirmationCount))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ObserverParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ObserverParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ObserverParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IsSupported {
		i--
		if m.IsSupported {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.MinObserverDelegation.Size()
		i -= size
		if _, err := m.MinObserverDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.BallotThreshold.Size()
		i -= size
		if _, err := m.BallotThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.Chain != nil {
		{
			size, err := m.Chain.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintParams(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ChainParamsList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ChainParams) > 0 {
		for _, e := range m.ChainParams {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *ChainParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ConfirmationCount != 0 {
		n += 1 + sovParams(uint64(m.ConfirmationCount))
	}
	if m.GasPriceTicker != 0 {
		n += 1 + sovParams(uint64(m.GasPriceTicker))
	}
	if m.InTxTicker != 0 {
		n += 1 + sovParams(uint64(m.InTxTicker))
	}
	if m.OutTxTicker != 0 {
		n += 1 + sovParams(uint64(m.OutTxTicker))
	}
	if m.WatchUtxoTicker != 0 {
		n += 1 + sovParams(uint64(m.WatchUtxoTicker))
	}
	l = len(m.ZetaTokenContractAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.ConnectorContractAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.Erc20CustodyContractAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.ChainId != 0 {
		n += 1 + sovParams(uint64(m.ChainId))
	}
	if m.OutboundTxScheduleInterval != 0 {
		n += 1 + sovParams(uint64(m.OutboundTxScheduleInterval))
	}
	if m.OutboundTxScheduleLookahead != 0 {
		n += 1 + sovParams(uint64(m.OutboundTxScheduleLookahead))
	}
	l = m.BallotThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinObserverDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.IsSupported {
		n += 3
	}
	return n
}

func (m *ObserverParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Chain != nil {
		l = m.Chain.Size()
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.BallotThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinObserverDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.IsSupported {
		n += 2
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ChainParamsList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: ChainParamsList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainParamsList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainParams", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainParams = append(m.ChainParams, &ChainParams{})
			if err := m.ChainParams[len(m.ChainParams)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *ChainParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: ChainParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ChainParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfirmationCount", wireType)
			}
			m.ConfirmationCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConfirmationCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasPriceTicker", wireType)
			}
			m.GasPriceTicker = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasPriceTicker |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field InTxTicker", wireType)
			}
			m.InTxTicker = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.InTxTicker |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutTxTicker", wireType)
			}
			m.OutTxTicker = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OutTxTicker |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WatchUtxoTicker", wireType)
			}
			m.WatchUtxoTicker = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WatchUtxoTicker |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZetaTokenContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ZetaTokenContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectorContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectorContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Erc20CustodyContractAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Erc20CustodyContractAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutboundTxScheduleInterval", wireType)
			}
			m.OutboundTxScheduleInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OutboundTxScheduleInterval |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OutboundTxScheduleLookahead", wireType)
			}
			m.OutboundTxScheduleLookahead = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OutboundTxScheduleLookahead |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BallotThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BallotThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinObserverDelegation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinObserverDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsSupported", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsSupported = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func (m *ObserverParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: ObserverParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ObserverParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Chain == nil {
				m.Chain = &chains.Chain{}
			}
			if err := m.Chain.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BallotThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BallotThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinObserverDelegation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinObserverDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsSupported", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsSupported = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
