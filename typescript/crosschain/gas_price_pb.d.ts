// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file crosschain/gas_price.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message zetachain.zetacore.crosschain.GasPrice
 */
export declare class GasPrice extends Message<GasPrice> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: string index = 2;
   */
  index: string;

  /**
   * @generated from field: int64 chain_id = 3;
   */
  chainId: bigint;

  /**
   * @generated from field: repeated string signers = 4;
   */
  signers: string[];

  /**
   * @generated from field: repeated uint64 block_nums = 5;
   */
  blockNums: bigint[];

  /**
   * @generated from field: repeated uint64 prices = 6;
   */
  prices: bigint[];

  /**
   * @generated from field: uint64 median_index = 7;
   */
  medianIndex: bigint;

  constructor(data?: PartialMessage<GasPrice>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.crosschain.GasPrice";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): GasPrice;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): GasPrice;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): GasPrice;

  static equals(a: GasPrice | PlainMessage<GasPrice> | undefined, b: GasPrice | PlainMessage<GasPrice> | undefined): boolean;
}

