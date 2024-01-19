package common

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Target fees in zeta for system transactions: 210000000000000azeta
const (
	PostGasPriceGasLimit            = 1_500_000
	AddTxHashToOutTxTrackerGasLimit = 200_000
	PostNonceGasLimit               = 200_000
	PostSendEVMGasLimit             = 4_000_000 // likely emit a lot of logs, so costly
	PostSendNonEVMGasLimit          = 1_000_000
	PostReceiveConfirmationGasLimit = 400_000
	PostBlameDataGasLimit           = 200_000
	DefaultGasLimit                 = 200_000
	PostProveOutboundTxGasLimit     = 400_000
	DefaultRetryCount               = 5
	ExtendedRetryCount              = 15
	DefaultRetryInterval            = 5

	StandardTXLimitFoSystemTx = "21000"
)

func GetAdjustTedGasPriceForTx(gasLimit int64, CurrentgasPrice sdkmath.Int) (sdk.Dec, error) {
	decStandardTx, err := sdk.NewDecFromStr(StandardTXLimitFoSystemTx)
	if err != nil {
		return sdk.Dec{}, err
	}
	decCurrentGasPrice := sdk.NewDecFromBigInt(CurrentgasPrice.BigInt())
	standardFees := decStandardTx.Mul(decCurrentGasPrice)

	adjustGasPrice := standardFees.Quo(sdk.NewDec(gasLimit))
	return adjustGasPrice, nil
}
