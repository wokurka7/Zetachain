package zetaclient

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/rs/zerolog"
)

const (
	satoshiPerBitcoin = 1e8
)

func getSatoshis(btc float64) (int64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	// BTC max amount is 21 mil and its at least 10^(-8) or one satoshi.
	switch {
	case math.IsNaN(btc):
		fallthrough
	case math.IsInf(btc, 1):
		fallthrough
	case math.IsInf(btc, -1):
		return 0, errors.New("invalid bitcoin amount")
	case btc > 21000000.0:
		return 0, errors.New("exceeded max bitcoin amount")
	case btc < 0.00000001:
		return 0, errors.New("cannot be less than 1 satoshi")
	}
	return round(btc * satoshiPerBitcoin), nil
}

func round(f float64) int64 {
	if f < 0 {
		return int64(f - 0.5)
	}
	return int64(f + 0.5)
}

func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

type DynamicTicker struct {
	name     string
	interval uint64
	impl     *time.Ticker
}

func NewDynamicTicker(name string, interval uint64) *DynamicTicker {
	return &DynamicTicker{
		name:     name,
		interval: interval,
		impl:     time.NewTicker(time.Duration(interval) * time.Second),
	}
}

func (t *DynamicTicker) C() <-chan time.Time {
	return t.impl.C
}

func (t *DynamicTicker) UpdateInterval(newInterval uint64, logger zerolog.Logger) {
	if newInterval > 0 && t.interval != newInterval {
		t.impl.Stop()
		t.interval = newInterval
		t.impl = time.NewTicker(time.Duration(t.interval) * time.Second)
		logger.Info().Msgf("%s ticker interval changed from %d to %d", t.name, t.interval, newInterval)
	}
}

func (t *DynamicTicker) Stop() {
	t.impl.Stop()
}

type DebugWriter struct {
	Sender        string
	SenderChain   int64
	Receiver      string
	ReceiverChain int64
	InTxHash      string
	InBlockHeight uint64
}

func WriteDebugDataToFile(logger zerolog.Logger, sender string, senderChain int64, receiver string, receiverChain int64, inTxHash string, inBlockHeight uint64) {
	const folder string = "debug_data"
	const filename string = "cctx_debug.json"
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Error().Msgf("Error accessing home directory")
	}

	folderPath := filepath.Join(home, folder)

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		logger.Error().Msgf("Error creating directory for debug data")
	}
	file := filepath.Join(home, folder, filename)
	file, err = filepath.Abs(file)
	if err != nil {
		logger.Error().Msgf("Error accessing file for debug data")
	}
	file = filepath.Clean(file)
	var debugData []DebugWriter
	input, err := os.ReadFile(file)
	if err != nil {
		logger.Info().Msgf("Error reading file for existing debug data")
	}
	if input != nil {
		err = json.Unmarshal(input, &debugData)
		if err != nil {
			logger.Error().Msgf("Error unmarshalling file for existing debug data")
		}
	}
	debugData = append(debugData, DebugWriter{
		Sender:        sender,
		SenderChain:   senderChain,
		Receiver:      receiver,
		ReceiverChain: receiverChain,
		InTxHash:      inTxHash,
		InBlockHeight: inBlockHeight,
	})
	jsonFile, _ := json.MarshalIndent(debugData, "", "  ")
	err = os.WriteFile(file, jsonFile, 0600)
	if err != nil {
		logger.Error().Msgf("Error writing file for debug data")
	}

}
