package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance

func NewParams() Params {
	defaultSlashAmount := sdk.ZeroInt()
	intSlashAmount, ok := sdkmath.NewIntFromString("100000000000000000")
	if ok {
		defaultSlashAmount = intSlashAmount
	}
	return Params{
		ValidatorEmissionPercentage: "00.50",
		ObserverEmissionPercentage:  "00.25",
		TssSignerEmissionPercentage: "00.25",
		ObserverSlashAmount:         defaultSlashAmount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPrefix(ParamValidatorEmissionPercentage), &p.ValidatorEmissionPercentage, validateValidatorEmissionPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamObserverEmissionPercentage), &p.ObserverEmissionPercentage, validateObserverEmissionPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamTssSignerEmissionPercentage), &p.TssSignerEmissionPercentage, validateTssEmissonPercentage),
		paramtypes.NewParamSetPair(KeyPrefix(ParamObserverSlashAmount), &p.ObserverSlashAmount, validateObserverSlashAmount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}

func validateObserverSlashAmount(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.LT(sdk.ZeroInt()) {
		return fmt.Errorf("slash amount cannot be less than 0")
	}
	return nil
}

func validateValidatorEmissionPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("validator emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("validator emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateObserverEmissionPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("observer emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("observer emission percentage cannot be less than 0 percent")
	}
	return nil
}

func validateTssEmissonPercentage(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	dec := sdk.MustNewDecFromStr(v)
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("tss emission percentage cannot be more than 100 percent")
	}
	if dec.LT(sdk.ZeroDec()) {
		return fmt.Errorf("tss emission percentage cannot be less than 0 percent")
	}
	return nil
}
