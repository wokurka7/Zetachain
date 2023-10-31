// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file crosschain/in_tx_hash_to_cctx.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message zetachain.zetacore.crosschain.InTxHashToCctx
 */
export declare class InTxHashToCctx extends Message<InTxHashToCctx> {
  /**
   * @generated from field: string in_tx_hash = 1;
   */
  inTxHash: string;

  /**
   * @generated from field: repeated string cctx_index = 2;
   */
  cctxIndex: string[];

  constructor(data?: PartialMessage<InTxHashToCctx>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.InTxHashToCctx";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): InTxHashToCctx;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): InTxHashToCctx;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): InTxHashToCctx;

  static equals(a: InTxHashToCctx | PlainMessage<InTxHashToCctx> | undefined, b: InTxHashToCctx | PlainMessage<InTxHashToCctx> | undefined): boolean;
}

