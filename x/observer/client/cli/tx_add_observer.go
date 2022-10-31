package cli

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

var _ = strconv.Itoa(0)

func CmdAddObserver() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-observer [observer-chain] [observation-type]",
		Short: "Broadcast message add-observer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argObserverChain := args[0]
			argObservationType := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			valAddress := sdk.ValAddress(clientCtx.GetFromAddress())
			operatorAddress, err := sdk.ValAddressFromBech32(valAddress.String())
			if err != nil {
				return err
			}
			msg := types.NewMsgAddObserver(
				clientCtx.GetFromAddress().String(),
				types.ParseStringToObserverChain(argObserverChain),
				types.ParseStringToObservationType(argObservationType),
				operatorAddress.String(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
