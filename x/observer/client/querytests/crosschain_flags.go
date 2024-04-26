package querytests

import (
	"fmt"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/observer/client/cli"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/status"
)

func (s *CliTestSuite) TestShowCrosschainFlags() {
	ctx := s.network.Validators[0].ClientCtx
	obj := s.observerState.CrosschainFlags
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  *types.CrosschainFlags
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	} {
		s.Run(tc.desc, func() {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowCrosschainFlags(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp types.QueryGetCrosschainFlagsResponse
				s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.CrosschainFlags)
				tc := tc
				s.Require().Equal(nullify.Fill(&tc.obj),
					nullify.Fill(&resp.CrosschainFlags),
				)
			}
		})
	}
}
