package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestKeeper_GetParams(t *testing.T) {
	tests := []struct {
		name    string
		params  emissionstypes.Params
		isPanic string
	}{
		{
			name: "Successfully set params",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "",
		},
		{
			name: "negative observer slashed amount",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(-10),
			},
			isPanic: "slash amount cannot be less than 0",
		},
		{
			name: "validator emission percentage too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "1.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "validator emission percentage cannot be more than 100 percent",
		},
		{
			name: "validator emission percentage too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "-1.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "validator emission percentage cannot be less than 0 percent",
		},
		{
			name: "observer percentage too low",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "-00.25",
				TssSignerEmissionPercentage: "00.25",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "observer emission percentage cannot be less than 0 percent",
		},
		{
			name: "observer percentage too high",
			params: emissionstypes.Params{

				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "150.25",
				TssSignerEmissionPercentage: "00.25",

				ObserverSlashAmount: sdkmath.NewInt(100000000000000000),
			},
			isPanic: "observer emission percentage cannot be more than 100 percent",
		},
		{
			name: "tss signer percentage too high",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "102.22",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "tss emission percentage cannot be more than 100 percent",
		},
		{
			name: "tss signer percentage too loo",
			params: emissionstypes.Params{
				ValidatorEmissionPercentage: "00.50",
				ObserverEmissionPercentage:  "00.25",
				TssSignerEmissionPercentage: "-102.22",
				ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
			},
			isPanic: "tss emission percentage cannot be less than 0 percent",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			assertPanic(t, func() {
				k.SetParams(ctx, tt.params)
			}, tt.isPanic)

			if tt.isPanic != "" {
				require.Equal(t, emissionstypes.DefaultParams(), k.GetParams(ctx))
			} else {
				require.Equal(t, tt.params, k.GetParams(ctx))
			}
		})
	}
}

func assertPanic(t *testing.T, f func(), errorLog string) {
	defer func() {
		r := recover()
		if r != nil {
			require.Contains(t, r, errorLog)
		}
	}()
	f()
}
