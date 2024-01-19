package ante

import (
	errorsmod "cosmossdk.io/errors"
)

const CodesSpace = "ante"

var (
	ErrTooManyMessages = errorsmod.Register(CodesSpace, 1, "too many messages")
	ErrNotASystemTx    = errorsmod.Register(CodesSpace, 2, "not a system tx")
)
