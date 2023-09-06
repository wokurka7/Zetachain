package zetaclient_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/zetaclient"
)

func TestWriteDebugDataToFile(t *testing.T) {
	const folder string = "debug_data"
	const filename string = "cctx_debug.json"
	home, err := os.UserHomeDir()
	assert.NoError(t, err)
	folderPath := filepath.Join(home, folder)
	file := filepath.Join(home, folder, filename)
	file, err = filepath.Abs(file)

	err = os.Remove(file)
	assert.NoError(t, err)

	for i := 0; i < 100; i++ {
		zetaclient.WriteDebugDataToFile(
			zerolog.New(os.Stdout).With().Str("module", "test").Logger(),
			"test",
			int64(i),
			"test",
			int64(i),
			"test",
			uint64(i))
	}

	err = os.MkdirAll(folderPath, os.ModePerm)
	assert.NoError(t, err)

	assert.NoError(t, err)
	file = filepath.Clean(file)
	var debugData []zetaclient.DebugWriter
	input, err := os.ReadFile(file)
	assert.NoError(t, err)
	err = json.Unmarshal(input, &debugData)
	assert.NoError(t, err)
	assert.Equal(t, 100, len(debugData))
}
