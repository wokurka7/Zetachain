package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgServer_UpdatePolicies(t *testing.T) {
	t.Run("can update policies", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		policies := sample.Policies()

		res, err := msgServer.UpdatePolicies(sdk.WrapSDKContext(ctx), &types.MsgUpdatePolicies{
			AuthorityAddress: sample.AccAddress(),
			Policies:         policies,
		})
		require.NotNil(t, res)
		require.NoError(t, err)

		// Check policy is set
		got, found := k.GetPolicies(ctx)
		require.True(t, found)
		require.Equal(t, policies, got)
	})
}