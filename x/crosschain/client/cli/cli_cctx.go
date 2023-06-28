package cli

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

var _ = strconv.Itoa(0)

func CmdListSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-cctx",
		Short: "list all CCTX",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllCctxRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.CctxAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowSend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-cctx [index]",
		Short: "shows a CCTX",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGetCctxRequest{
				Index: args[0],
			}

			res, err := queryClient.Cctx(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// Transaction CLI /////////////////////////
//zetacored tx zetacore cctx-voter 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 0x96B05C238b99768F349135de0653b687f9c13fEE ETH 1000000000000000000 0 message hash 100 --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1

func CmdCCTXInboundVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inbound-voter [sender] [senderChainID] [txOrigin] [receiver] [receiverChainID] [amount]  [message] [inTxHash] [inBlockHeight] [coinType] [asset]",
		Short: "Broadcast message sendVoter",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSender := (args[0])
			argsSenderChain, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			argsTxOrigin := (args[2])
			argsReceiver := (args[3])
			argsReceiverChain, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			amount := math.NewUintFromString(args[5])
			argsMessage := (args[6])
			argsInTxHash := (args[7])
			argsInBlockHeight, err := strconv.ParseInt(args[8], 10, 64)
			argsCoinType := common.CoinType(common.CoinType_value[args[9]])
			argsAsset := args[10]
			if err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendVoter(clientCtx.GetFromAddress().String(), (argsSender), int64((argsSenderChain)), (argsTxOrigin), (argsReceiver), int64((argsReceiverChain)), amount, (argsMessage), (argsInTxHash), uint64(argsInBlockHeight), 250_000, argsCoinType, argsAsset)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdCCTXOutboundVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outbound-voter [sendHash] [outTxHash] [outBlockHeight] [ZetaMinted] [Status] [chain] [outTXNonce] [coinType]",
		Short: "Broadcast message receiveConfirmation",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSendHash := (args[0])
			argsOutTxHash := (args[1])
			argsOutBlockHeight, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				return err
			}
			argsMMint := (args[3])
			var status common.ReceiveStatus
			if args[4] == "0" {
				status = common.ReceiveStatus_Success
			} else if args[4] == "1" {
				status = common.ReceiveStatus_Failed
			} else {
				return fmt.Errorf("wrong status")
			}
			chain, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}
			outTxNonce, err := strconv.ParseInt(args[6], 10, 64)
			if err != nil {
				return err
			}
			argsCoinType := common.CoinType(common.CoinType_value[args[7]])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgReceiveConfirmation(clientCtx.GetFromAddress().String(), argsSendHash, argsOutTxHash, uint64(argsOutBlockHeight), math.NewUintFromString(argsMMint), status, int64(chain), uint64(outTxNonce), argsCoinType)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// Transaction CLI /////////////////////////
//zetacored tx zetacore whitelist-voter ETH whitelist 0x96B05C238b99768F349135de0653b687f9c13fEE --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1
//zetacored tx zetacore whitelist-voter ETH unwhitelist 0x96B05C238b99768F349135de0653b687f9c13fEE --from=zeta --keyring-backend=test --yes --chain-id=localnet_101-1

func CmdCCTXWhitelistVoter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whitelist-voter [receiverChainID] [message] [asset]",
		Short: "Broadcast message sendVoter",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsReceiverChain, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			argsMessage := args[1]
			argsAsset := args[2]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSendVoter(clientCtx.GetFromAddress().String(), "", common.ZetaChain().ChainId, "", "", int64((argsReceiverChain)), math.ZeroUint(), (argsMessage), "", 0, 250_000, common.CoinType_ERC20, argsAsset)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
