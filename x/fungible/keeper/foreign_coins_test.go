package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"strconv"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNForeignCoins(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ForeignCoins {
	items := make([]types.ForeignCoins, n)
	for i := range items {
		items[i].Zrc20ContractAddress = strconv.Itoa(i)

		keeper.SetForeignCoins(ctx, items[i])
	}
	return items
}
