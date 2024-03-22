package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
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
	err = AddSdkState(appState, importAppState, cdc)
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
	//if err := AddCrossChainState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	//if err := AddObserverState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	//if err := AddEmissionsState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	//if err := AddFungibleState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	//if err := AddAuthorityState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	return nil
}

func AddSdkState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	//if err := AddAuthState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	//if err := AddStakingState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	if err := AddEvmState(appState, importAppState, cdc); err != nil {
		return err
	}
	//if err := AddDistributionState(appState, importAppState, cdc); err != nil {
	//	return err
	//}
	return nil
}

func AddAuthState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedCrossChainGenState := authtypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedCrossChainStateBz, err := cdc.MarshalJSON(&importedCrossChainGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal zetacrosschain genesis state: %w", err)
	}
	appState[authtypes.ModuleName] = importedCrossChainStateBz
	return nil
}

func AddStakingState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	importedStakingGenState := stakingtypes.GetGenesisStateFromAppState(cdc, importAppState)
	importedStakingStateBz, err := cdc.MarshalJSON(importedStakingGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal zetacrosschain genesis state: %w", err)
	}
	appState[stakingtypes.ModuleName] = importedStakingStateBz
	return nil
}

func AddDistributionState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	var importedDistributionGenState distributiontypes.GenesisState
	fmt.Println("importAppState: ", unsafe.Sizeof(importAppState))
	if importAppState[distributiontypes.ModuleName] == nil {
		panic("distribution module not found in import file")
	}

	err := cdc.UnmarshalJSON(appState[distributiontypes.ModuleName], &importedDistributionGenState)
	if err != nil {
		return fmt.Errorf("failed to unmarshal distribution genesis state: %w", err)
	}
	fmt.Println("Number of delegations: ", len(importedDistributionGenState.DelegatorStartingInfos))
	fmt.Println("genesis: ", importedDistributionGenState.String())
	importedDistributionStateBz, err := cdc.MarshalJSON(&importedDistributionGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal distribution genesis state: %w", err)
	}
	appState[distributiontypes.ModuleName] = importedDistributionStateBz
	return nil
}

func AddEvmState(appState map[string]json.RawMessage, importAppState map[string]json.RawMessage, cdc codec.Codec) error {
	var importedEvmGenState evmtypes.GenesisState
	if importAppState[evmtypes.ModuleName] != nil {
		err := cdc.UnmarshalJSON(appState[evmtypes.ModuleName], &importedEvmGenState)
		if err != nil {
			return fmt.Errorf("failed to unmarshal evm genesis state: %w", err)
		}
	}

	err := codectypes.UnpackInterfaces(importedEvmGenState, cdc)
	if err != nil {
		return fmt.Errorf("failed to authz grants into upackeder: %w", err)
	}
	fmt.Println("Number of accounts: ", len(importedEvmGenState.Accounts))
	importedEvmStateBz, err := cdc.MarshalJSON(&importedEvmGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal evm genesis state: %w", err)
	}
	appState[evmtypes.ModuleName] = importedEvmStateBz
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
