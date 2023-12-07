// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file observer/tss_funds_migrator.proto (package zetachain.zetacore.observer, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message zetachain.zetacore.observer.TssFundMigratorInfo
 */
export declare class TssFundMigratorInfo extends Message<TssFundMigratorInfo> {
  /**
   * @generated from field: int64 chain_id = 1;
   */
  chainId: bigint;

  /**
   * @generated from field: string migration_cctx_index = 2;
   */
  migrationCctxIndex: string;

  constructor(data?: PartialMessage<TssFundMigratorInfo>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "zetachain.zetacore.observer.TssFundMigratorInfo";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): TssFundMigratorInfo;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): TssFundMigratorInfo;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): TssFundMigratorInfo;

  static equals(a: TssFundMigratorInfo | PlainMessage<TssFundMigratorInfo> | undefined, b: TssFundMigratorInfo | PlainMessage<TssFundMigratorInfo> | undefined): boolean;
}

