package main

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nanmu42/etherscan-api"
	"github.com/rs/zerolog"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	ZETA_RPC_NODE_URL := "46.4.15.110:9090"
	grpcConn, err := grpc.Dial(ZETA_RPC_NODE_URL, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer grpcConn.Close()
	var wg sync.WaitGroup
	crosschainClient := crosschaintypes.NewQueryClient(grpcConn)
	wg.Add(1)
	go func() {
		defer wg.Done()
		//BitcoinCrossCheck(crosschainClient, 1000)
	}()
	goerliClient, err := ethclient.Dial("https://summer-red-dust.ethereum-goerli.discover.quiknode.pro/d59692c13082ce6b8290e01db1324ac1e27ee54a/")
	if err != nil {
		panic(err)
	}
	//goerliEtherscanAPIURL := "https://api-goerli.etherscan.io/api"
	goerliEtherscanClient := etherscan.NewCustomized(etherscan.Customization{
		Timeout: 15 * time.Second,
		Key:     "E1DG89WEBS2ZESXQNWUDKJJIIYQEX4HIHP",
		BaseURL: "https://api-goerli.etherscan.io/api",
		Verbose: false,
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		EthishCrossCheckTSSInbound(crosschainClient, goerliClient, goerliEtherscanClient)
	}()
	wg.Wait()
}

func BitcoinCrossCheck(crosschainClient crosschaintypes.QueryClient, numBlocksToCheck int64) {
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
	startBN := bn - numBlocksToCheck
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

func EthishCrossCheckTSSInbound(crosschainClient crosschaintypes.QueryClient, client *ethclient.Client, etherscanClient *etherscan.Client) {
	tssResp, err := crosschainClient.GetTssAddress(context.Background(), &crosschaintypes.QueryGetTssAddressRequest{})
	if err != nil {
		panic(err)
	}
	tssAddress := ethcommon.HexToAddress(tssResp.Eth)
	if tssAddress == (ethcommon.Address{}) {
		panic("tss address not found")
	}
	bn, err := client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	startBN := int(bn - 75)
	endBN := int(bn - 15)

	kiaCnt := 0
	inTxCnt := 0
	for i := startBN; i < endBN; i++ {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(i)))
		if err != nil {
			panic(err)
		}
		//ts := time.Unix(int64(block.Time()), 0)
		//fmt.Printf("ETH Block %d: [%d/%d], %s\n", i, i-startBN, endBN-startBN, ts.String())
		for _, tx := range block.Transactions() {
			if tx.To() != nil && *tx.To() == tssAddress {
				inTxCnt++
				//fmt.Printf("  InTx: %s\n", tx.Hash().String())
				res, err := crosschainClient.InTxHashToCctx(context.Background(), &crosschaintypes.QueryGetInTxHashToCctxRequest{
					InTxHash: tx.Hash().String(),
				})
				if err != nil {
					st, ok := status.FromError(err)
					if ok && st.Code() == codes.NotFound {
						//fmt.Printf("Error is of type NotFound: %s\n", st.Message())
						fmt.Printf("  #### XCHECK FAIL: CCTX NOT REGISTERED FOR INCOMING TX %s\n", tx.Hash().String())
						kiaCnt++
					} else {
						fmt.Println("  Error is not of type NotFound!")
						panic(err)
					}
				} else {
					_, err := crosschainClient.Cctx(context.Background(), &crosschaintypes.QueryGetCctxRequest{
						Index: res.InTxHashToCctx.CctxIndex[0],
					})
					if err != nil {
						panic(err)
					}
					//fmt.Printf("  OK: cctx status %s\n", cctx.CrossChainTx.CctxStatus.Status.String())
				}
			}
		}
	}

	fmt.Printf("#### XCHECK COMPLETE: %d/%d CCTX NOT REGISTERED in the past %d blocks\n", kiaCnt, inTxCnt, endBN-startBN)
}
