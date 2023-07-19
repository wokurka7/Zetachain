// @generated by protoc-gen-es v1.3.0
// @generated from file observer/tx.proto (package zetachain.zetacore.observer, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { CoreParams } from "./params_pb.js";
import type { ObservationType } from "./observer_pb.js";

/**
 * message MsgSetSupportedChains {
 *  string creator = 1;
 *  int64 chain_id=2;
 *  common.ChainName ChainName=3;
 * }
 *
 * @generated from message zetachain.zetacore.observer.MsgUpdateCoreParams
 */
export declare class MsgUpdateCoreParams extends Message<MsgUpdateCoreParams> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: zetachain.zetacore.observer.CoreParams coreParams = 2;
   */
  coreParams?: CoreParams;

  constructor(data?: PartialMessage<MsgUpdateCoreParams>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateCoreParams";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateCoreParams;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateCoreParams;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateCoreParams;

  static equals(a: MsgUpdateCoreParams | PlainMessage<MsgUpdateCoreParams> | undefined, b: MsgUpdateCoreParams | PlainMessage<MsgUpdateCoreParams> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgUpdateCoreParamsResponse
 */
export declare class MsgUpdateCoreParamsResponse extends Message<MsgUpdateCoreParamsResponse> {
  constructor(data?: PartialMessage<MsgUpdateCoreParamsResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgUpdateCoreParamsResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgUpdateCoreParamsResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgUpdateCoreParamsResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgUpdateCoreParamsResponse;

  static equals(a: MsgUpdateCoreParamsResponse | PlainMessage<MsgUpdateCoreParamsResponse> | undefined, b: MsgUpdateCoreParamsResponse | PlainMessage<MsgUpdateCoreParamsResponse> | undefined): boolean;
}

/**
 * this line is used by starport scaffolding # proto/tx/message
 *
 * message MsgSetSupportedChainsResponse{
 * }
 *
 * @generated from message zetachain.zetacore.observer.MsgAddObserver
 */
export declare class MsgAddObserver extends Message<MsgAddObserver> {
  /**
   * @generated from field: string creator = 1;
   */
  creator: string;

  /**
   * @generated from field: int64 chain_id = 2;
   */
  chainId: bigint;

  /**
   * @generated from field: zetachain.zetacore.observer.ObservationType observationType = 3;
   */
  observationType: ObservationType;

  constructor(data?: PartialMessage<MsgAddObserver>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddObserver";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddObserver;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddObserver;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddObserver;

  static equals(a: MsgAddObserver | PlainMessage<MsgAddObserver> | undefined, b: MsgAddObserver | PlainMessage<MsgAddObserver> | undefined): boolean;
}

/**
 * @generated from message zetachain.zetacore.observer.MsgAddObserverResponse
 */
export declare class MsgAddObserverResponse extends Message<MsgAddObserverResponse> {
  constructor(data?: PartialMessage<MsgAddObserverResponse>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.MsgAddObserverResponse";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): MsgAddObserverResponse;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): MsgAddObserverResponse;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): MsgAddObserverResponse;

  static equals(a: MsgAddObserverResponse | PlainMessage<MsgAddObserverResponse> | undefined, b: MsgAddObserverResponse | PlainMessage<MsgAddObserverResponse> | undefined): boolean;
}

