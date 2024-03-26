package cli

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CmdBurnTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-tokens [chain-id] [amount] [burn-addr]",
		Short: `Burn zeta tokens on the specified chain using the connector contract`,
		Args:  cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainID, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			amount := sdkmath.NewUintFromString(args[1])
			msg := types.NewMsgReduceZetaSupply(clientCtx.GetFromAddress().String(), chainID, amount, args[2])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
