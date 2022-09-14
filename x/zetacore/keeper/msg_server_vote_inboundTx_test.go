package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CalculateGassFee(t *testing.T) {

	tt := []struct {
		name        string
		gasPrice    sdk.Uint
		gasLimit    sdk.Uint
		rate        sdk.Uint
		expectedFee sdk.Uint
	}{
		{
			name:        "Test Price1",
			gasPrice:    sdk.NewUintFromString("20000000000"),
			gasLimit:    sdk.NewUintFromString("90000"),
			rate:        sdk.NewUintFromString("1000000000000000000"),
			expectedFee: sdk.NewUintFromString("1001800000000000000"),
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedFee, CalculateFee(test.gasPrice, test.gasLimit, test.rate))
		})
	}
}
