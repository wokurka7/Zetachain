package querytests

import (
	"fmt"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/zeta-chain/zetacore/x/observer/client/cli"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/status"
)

func (s *CliTestSuite) TestShowTSS() {
	ctx := s.network.Validators[0].ClientCtx
	obj := s.observerState.Tss
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		args []string
		err  error
		obj  *types.TSS
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
			err:  nil,
		},
	} {
		tc := tc
		s.Run(tc.desc, func() {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowTSS(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				s.Require().True(ok)
				s.Require().ErrorIs(stat.Err(), tc.err)
			} else {
				s.Require().NoError(err)
				var resp types.QueryGetTSSResponse
				s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().NotNil(resp.TSS)
				s.Require().Equal(*tc.obj, resp.TSS)
			}
		})
	}
}
