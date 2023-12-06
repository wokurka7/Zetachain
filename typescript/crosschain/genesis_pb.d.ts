// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file crosschain/genesis.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { Params } from "./params_pb.js";
import type { OutTxTracker } from "./out_tx_tracker_pb.js";
import type { GasPrice } from "./gas_price_pb.js";
import type { ChainNonces } from "./chain_nonces_pb.js";
import type { CrossChainTx, ZetaAccounting } from "./cross_chain_tx_pb.js";
import type { LastBlockHeight } from "./last_block_height_pb.js";
import type { InTxHashToCctx } from "./in_tx_hash_to_cctx_pb.js";
import type { InTxTracker } from "./in_tx_tracker_pb.js";

/**
 * GenesisState defines the metacore module's genesis state.
 *
 * @generated from message zetachain.zetacore.crosschain.GenesisState
 */
export declare class GenesisState extends Message<GenesisState> {
  /**
   * @generated from field: zetachain.zetacore.crosschain.Params params = 1;
   */
  params?: Params;

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.OutTxTracker outTxTrackerList = 2;
   */
  outTxTrackerList: OutTxTracker[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.GasPrice gasPriceList = 5;
   */
  gasPriceList: GasPrice[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.ChainNonces chainNoncesList = 6;
   */
  chainNoncesList: ChainNonces[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.CrossChainTx CrossChainTxs = 7;
   */
  CrossChainTxs: CrossChainTx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.LastBlockHeight lastBlockHeightList = 8;
   */
  lastBlockHeightList: LastBlockHeight[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxHashToCctx inTxHashToCctxList = 9;
   */
  inTxHashToCctxList: InTxHashToCctx[];

  /**
   * @generated from field: repeated zetachain.zetacore.crosschain.InTxTracker in_tx_tracker_list = 11;
   */
  inTxTrackerList: InTxTracker[];

  /**
   * @generated from field: zetachain.zetacore.crosschain.ZetaAccounting zeta_accounting = 12;
   */
  zetaAccounting?: ZetaAccounting;

  constructor(data?: PartialMessage<GenesisState>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.GenesisState";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GenesisState;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GenesisState;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GenesisState;

  static equals(a: GenesisState | PlainMessage<GenesisState> | undefined, b: GenesisState | PlainMessage<GenesisState> | undefined): boolean;
}

