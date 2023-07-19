// @generated by protoc-gen-es v1.3.0
// @generated from file fungible/genesis.proto (package zetachain.zetacore.fungible, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { proto3 } from "@bufbuild/protobuf";
import { Params } from "./params_pb.js";
import { ForeignCoins } from "./foreign_coins_pb.js";
import { SystemContract } from "./system_contract_pb.js";

/**
 * GenesisState defines the fungible module's genesis state.
 *
 * @generated from message zetachain.zetacore.fungible.GenesisState
 */
export const GenesisState = proto3.makeMessageType(
  "zetachain.zetacore.fungible.GenesisState",
  () => [
    { no: 1, name: "params", kind: "message", T: Params },
    { no: 2, name: "foreignCoinsList", kind: "message", T: ForeignCoins, repeated: true },
    { no: 3, name: "systemContract", kind: "message", T: SystemContract },
  ],
);

