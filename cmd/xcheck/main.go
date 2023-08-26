package main

import (
	"context"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func main() {
	ZETA_RPC_NODE_URL := "46.4.15.110:9090"
	grpcConn, err := grpc.Dial(ZETA_RPC_NODE_URL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()

	crosschainClient := crosschaintypes.NewQueryClient(grpcConn)
	res, err := crosschainClient.InTxHashToCctx(context.Background(), &crosschaintypes.QueryGetInTxHashToCctxRequest{
		InTxHash: "8fd1c39d2f2f6630a92353164d7a230aed0e69c1efe8593f6eda872e961691fb",
	})
	st, ok := status.FromError(err)
	if ok && st.Code() == codes.NotFound {
		fmt.Printf("Error is of type NotFound: %s\n", st.Message())
	} else {
		fmt.Println("Error is not of type NotFound!")
	}
	_ = res

	tssResp, err := crosschainClient.GetTssAddress(context.Background(), &crosschaintypes.QueryGetTssAddressRequest{})
	if err != nil {
		panic(err)
	}
	btcTSSAddressString := tssResp.Btc
	chainParams := chaincfg.TestNet3Params
	config.BitconNetParams = &chainParams
	btcTSSAddress, err := btcutil.DecodeAddress(btcTSSAddressString, &chainParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("TSS BTC Address: %s\n", btcTSSAddress.EncodeAddress())

	connCfg := &rpcclient.ConnConfig{
		Host:         "bitcoin-rpc.athens.zetachain.com",
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       "testnet3",
	}
	btcClient, err := rpcclient.New(connCfg, nil)
	bn, err := btcClient.GetBlockCount()
	if err != nil {
		panic(err)
	}
	startBN := bn - 1000
	endBN := bn
	kiaCnt := 0
	inTxCnt := 0
	for i := startBN; i < endBN; i++ {
		hash, err := btcClient.GetBlockHash(i)
		if err != nil {
			panic(err)
		}
		block, err := btcClient.GetBlockVerboseTx(hash)
		if err != nil {
			panic(err)
		}
		ts := time.Unix(block.Time, 0)
		fmt.Printf("BTC Block %d: [%d/%d], %s\n", i, i-startBN, endBN-startBN, ts.String())

		logger := zerolog.Logger{}
		inTxs := zetaclient.FilterAndParseIncomingTx(block.Tx, uint64(block.Height), btcTSSAddress.EncodeAddress(), &logger)
		for _, inTx := range inTxs {
			inTxCnt++
			fmt.Printf("  InTx: %s\n", inTx.TxHash)
			res, err := crosschainClient.InTxHashToCctx(context.Background(), &crosschaintypes.QueryGetInTxHashToCctxRequest{
				InTxHash: inTx.TxHash,
			})
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.NotFound {
					//fmt.Printf("Error is of type NotFound: %s\n", st.Message())
					fmt.Printf("  #### XCHECK FAIL: CCTX NOT REGISTERED FOR INCOMING TX %s\n", inTx.TxHash)
					kiaCnt++
				} else {
					fmt.Println("  Error is not of type NotFound!")
					panic(err)
				}
			} else {
				cctx, err := crosschainClient.Cctx(context.Background(), &crosschaintypes.QueryGetCctxRequest{
					Index: res.InTxHashToCctx.CctxIndex[0],
				})
				if err != nil {
					panic(err)
				}
				fmt.Printf("  OK: cctx status %s\n", cctx.CrossChainTx.CctxStatus.Status.String())
			}
			_ = res
		}
	}

	fmt.Printf("#### XCHECK COMPLETE: %d/%d CCTX NOT REGISTERED in the past %d blocks\n", kiaCnt, inTxCnt, endBN-startBN)
}
