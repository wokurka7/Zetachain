package common

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
)

func GetAdjustTedGasPriceForTx(gasLimit int, expectedFee) {

}
