package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) GetReservesFactor(ctx sdk.Context) sdk.Dec {
	reserveAmount := k.GetBankKeeper().GetBalance(ctx, types.EmissionsModuleAddress, config.BaseDenom)
	return sdk.NewDecFromInt(reserveAmount.Amount)
}

func CalculateFixedValidatorRewards(avgBlockTimeString string) (sdk.Dec, error) {
	azetaAmountTotalRewards, err := common.GetAzetaDecFromAmountInZeta(types.BlockRewardsInZeta)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	avgBlockTime, err := sdk.NewDecFromStr(avgBlockTimeString)
	if err != nil {
		return sdk.ZeroDec(), err
	}
	numberOfBlocksInAMonth := sdk.NewDec(types.SecsInMonth).Quo(avgBlockTime)
	numberOfBlocksTotal := numberOfBlocksInAMonth.Mul(sdk.NewDec(12)).Mul(sdk.NewDec(types.EmissionScheduledYears))
	constantRewardPerBlock := azetaAmountTotalRewards.Quo(numberOfBlocksTotal)
	return constantRewardPerBlock, nil
}
