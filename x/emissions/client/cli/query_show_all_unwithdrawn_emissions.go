package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func CmdListUnwithdrawnEmissions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-unwithdrawn-emissions",
		Short: "Query list-unwithdrawn-emissions",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.ListUnWithdrawnEmissions(cmd.Context(), &types.QueryListUnWithdrawnEmissionsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)

		},
	}

	return cmd

}
