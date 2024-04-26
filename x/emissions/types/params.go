package types

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"gopkg.in/yaml.v2"
)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MaxBondFactor:               "1.25",
		MinBondFactor:               "0.75",
		AvgBlockTime:                "6.00",
		TargetBondRatio:             "00.67",
		ValidatorEmissionPercentage: "00.50",
		ObserverEmissionPercentage:  "00.25",
		TssSignerEmissionPercentage: "00.25",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         sdkmath.NewInt(100000000000000000),
		BallotMaturityBlocks:        100,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMaxBondFactor(p.MaxBondFactor); err != nil {
		return err
	}
	if err := validateMinBondFactor(p.MinBondFactor); err != nil {
		return err
	}
	if err := validateAvgBlockTime(p.AvgBlockTime); err != nil {
		return err
	}
	if err := validateTargetBondRatio(p.TargetBondRatio); err != nil {
		return err
	}
	if err := validateValidatorEmissionPercentage(p.ValidatorEmissionPercentage); err != nil {
		return err
	}
	if err := validateObserverEmissionPercentage(p.ObserverEmissionPercentage); err != nil {
		return err
	}
	if err := validateTssEmissionPercentage(p.TssSignerEmissionPercentage); err != nil {
		return err
	}
	if err := validateBallotMaturityBlocks(p.BallotMaturityBlocks); err != nil {
		return err
	}
	return validateObserverSlashAmount(p.ObserverSlashAmount)
}

func (p Params) GetBondFactor(currentBondedRatio sdk.Dec) sdk.Dec {
	targetBondRatio := sdk.MustNewDecFromStr(p.TargetBondRatio)
	maxBondFactor := sdk.MustNewDecFromStr(p.MaxBondFactor)
	minBondFactor := sdk.MustNewDecFromStr(p.MinBondFactor)

	// Bond factor ranges between minBondFactor (0.75) to maxBondFactor (1.25)
	if currentBondedRatio.IsZero() {
		return sdk.ZeroDec()
	}
	bondFactor := targetBondRatio.Quo(currentBondedRatio)
	if bondFactor.GT(maxBondFactor) {
		return maxBondFactor
	}
	if bondFactor.LT(minBondFactor) {
		return minBondFactor
	}
	return bondFactor
}

func (p Params) GetDurationFactor(blockHeight int64) sdk.Dec {
	avgBlockTime := sdk.MustNewDecFromStr(p.AvgBlockTime)
	NumberOfBlocksInAMonth := sdk.NewDec(SecsInMonth).Quo(avgBlockTime)
	monthFactor := sdk.NewDec(blockHeight).Quo(NumberOfBlocksInAMonth)
	logValueDec := sdk.MustNewDecFromStr(p.DurationFactorConstant)
	// month * log(1 + 0.02 / 12)
	fractionNumerator := monthFactor.Mul(logValueDec)
	// (month * log(1 + 0.02 / 12) ) + 1
	fractionDenominator := fractionNumerator.Add(sdk.OneDec())

	// (month * log(1 + 0.02 / 12)) / (month * log(1 + 0.02 / 12) ) + 1
	if fractionDenominator.IsZero() {
		return sdk.OneDec()
	}
	if fractionNumerator.IsZero() {
		return sdk.ZeroDec()
	}
	return fractionNumerator.Quo(fractionDenominator)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		return ""
	}
	return string(out)
}

func validateDurationFactorConstant(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateMaxBondFactor(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.GT(sdk.MustNewDecFromStr("1.25")) {
		return fmt.Errorf("max bond factor cannot be higher that 1.25")
	}
	return nil
}

func validateMinBondFactor(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.LT(sdk.MustNewDecFromStr("0.75")) {
		return fmt.Errorf("min bond factor cannot be lower that 0.75")
	}
	return nil
}

func validateAvgBlockTime(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	blocktime, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("invalid block time: %T", i)
	}
	if blocktime <= 0 {
		return fmt.Errorf("block time cannot be less than or equal to 0")
	}
	return nil
}

func validateTargetBondRatio(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	decMaxBond := sdk.MustNewDecFromStr(v)
	if decMaxBond.GT(sdk.OneDec()) {
		return fmt.Errorf("target bond ratio cannot be more than 100 percent")
	}
	if decMaxBond.LT(sdk.ZeroDec()) {
		return fmt.Errorf("target bond ratio cannot be less than 0 percent")
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

func validateTssEmissionPercentage(i interface{}) error {
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

func validateBallotMaturityBlocks(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 0 {
		return fmt.Errorf("ballot maturity types must be gte 0")
	}

	return nil
}
