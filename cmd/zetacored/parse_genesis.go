package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/types"
	"github.com/zeta-chain/zetacore/app"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func CmdParseGenesisFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse-genesis-file [import-genesis-file] [optional-genesis-file]",
		Short: "Parse the genesis file",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec
			genesisFilePAth := filepath.Join(app.DefaultNodeHome, "config", "genesis.json")
			if len(args) == 2 {
				genesisFilePAth = args[1]
			}
			genesis, err := GetGenDoc(genesisFilePAth)
			importData, err := GetGenDoc(args[0])

			err = ImportDataIntoFile(genesis, importData, cdc)
			if err != nil {
				return err
			}

			err = genutil.ExportGenesisFile(genesis, genesisFilePAth)
			if err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func GetGenDoc(fp string) (*types.GenesisDoc, error) {
	path, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}
	jsonBlob, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	genData, err := types.GenesisDocFromJSON(jsonBlob)
	if err != nil {
		return nil, err
	}
	return genData, nil
}

func ImportDataIntoFile(genDoc *types.GenesisDoc, importFile *types.GenesisDoc, cdc codec.Codec) error {
	appState, err := genutiltypes.GenesisStateFromGenDoc(*genDoc)
	if err != nil {
		return err
	}
	importAppState, err := genutiltypes.GenesisStateFromGenDoc(*importFile)
	if err != nil {
		return err
	}
	err = AddZetaState(appState, importAppState, cdc)
	if err != nil {
		return err
	}
	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}
	genDoc.AppState = appStateJSON
	return nil
}

func AddZetaState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	if err := AddCrossChainState(appState, importAppState, cdc); err != nil {
		return err
	}
	if err := AddObserverState(appState, importAppState, cdc); err != nil {
		return err
	}
	if err := AddEmissionsState(appState, importAppState, cdc); err != nil {
		return err
	}
	if err := AddFungibleState(appState, importAppState, cdc); err != nil {
		return err
	}
	if err := AddAuthorityState(appState, importAppState, cdc); err != nil {
		return err
	}
	return nil
}

func AddCrossChainState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedCrossChainGenState := crosschaintypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedCrossChainStateBz, err := json.Marshal(importedCrossChainGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal zetacrosschain genesis state: %w", err)
	}
	appState[crosschaintypes.ModuleName] = importedCrossChainStateBz
	return nil
}

func AddObserverState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedObserverGenState := observertypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedObserverStateBz, err := json.Marshal(importedObserverGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal observer genesis state: %w", err)
	}
	appState[observertypes.ModuleName] = importedObserverStateBz
	return nil
}

func AddEmissionsState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedEmissionsGenState := emissionstypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedEmissionsStateBz, err := json.Marshal(importedEmissionsGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal emissions genesis state: %w", err)
	}
	appState[emissionstypes.ModuleName] = importedEmissionsStateBz
	return nil
}

func AddFungibleState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedFungibleGenState := fungibletypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedFungibleStateBz, err := json.Marshal(importedFungibleGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal fungible genesis state: %w", err)
	}
	appState[fungibletypes.ModuleName] = importedFungibleStateBz
	return nil
}

func AddAuthorityState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedAuthorityGenState := authoritytypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedAuthorityStateBz, err := json.Marshal(importedAuthorityGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal authority genesis state: %w", err)
	}
	appState[authoritytypes.ModuleName] = importedAuthorityStateBz
	return nil
}
