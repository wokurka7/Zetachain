// @generated by protoc-gen-es v1.3.0
// @generated from file crosschain/last_block_height.proto (package zetachain.zetacore.crosschain, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message zetachain.zetacore.crosschain.LastBlockHeight
 */
export const LastBlockHeight = proto3.makeMessageType(
  "zetachain.zetacore.crosschain.LastBlockHeight",
  () => [
    { no: 1, name: "creator", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 2, name: "index", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 3, name: "chain", kind: "scalar", T: 9 /* ScalarType.STRING */ },
    { no: 4, name: "lastSendHeight", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
    { no: 5, name: "lastReceiveHeight", kind: "scalar", T: 4 /* ScalarType.UINT64 */ },
  ],
);

