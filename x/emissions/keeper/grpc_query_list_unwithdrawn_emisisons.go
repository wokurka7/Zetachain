package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
	"golang.org/x/net/context"
)

func (k Keeper) ListUnWithdrawnEmissions(goCtx context.Context, _ *types.QueryListUnWithdrawnEmissionsRequest) (*types.QueryListUnWithdrawnEmissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	emissions := k.GetAllWithdrawableEmission(ctx)
	total := sdkmath.ZeroInt()
	for _, emission := range emissions {
		total = total.Add(emission.Amount)
	}
	return &types.QueryListUnWithdrawnEmissionsResponse{Amount: total.String()}, nil
}
