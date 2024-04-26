package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	lightclientmocks "github.com/zeta-chain/zetacore/testutil/keeper/mocks/lightclient"
	"github.com/zeta-chain/zetacore/x/lightclient/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// LightclientMockOptions represents options for instantiating a lightclient keeper with mocks
type LightclientMockOptions struct {
	UseAuthorityMock bool
}

var (
	LightclientMocksAll = LightclientMockOptions{
		UseAuthorityMock: true,
	}
	LightclientNoMocks = LightclientMockOptions{}
)

func initLightclientKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
	authorityKeeper types.AuthorityKeeper,
) keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	return keeper.NewKeeper(cdc, storeKey, memKey, authorityKeeper)
}

// LightclientKeeperWithMocks instantiates a lightclient keeper for testing purposes with the option to mock specific keepers
func LightclientKeeperWithMocks(t testing.TB, mockOptions LightclientMockOptions) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

	authorityKeeperTmp := initAuthorityKeeper(cdc, db, stateStore)

	// Create regular keepers
	sdkKeepers := NewSDKKeepers(cdc, db, stateStore)

	// Create the observer keeper
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	ctx := NewContext(stateStore)

	// Initialize modules genesis
	sdkKeepers.InitGenesis(ctx)

	// Add a proposer to the context
	ctx = sdkKeepers.InitBlockProposer(t, ctx)

	// Initialize mocks for mocked keepers
	var authorityKeeper types.AuthorityKeeper = authorityKeeperTmp
	if mockOptions.UseAuthorityMock {
		authorityKeeper = lightclientmocks.NewLightclientAuthorityKeeper(t)
	}

	k := keeper.NewKeeper(cdc, storeKey, memStoreKey, authorityKeeper)

	return &k, ctx, sdkKeepers, ZetaKeepers{
		AuthorityKeeper: &authorityKeeperTmp,
	}
}

// LightclientKeeper instantiates an lightclient keeper for testing purposes
func LightclientKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, SDKKeepers, ZetaKeepers) {
	return LightclientKeeperWithMocks(t, LightclientNoMocks)
}

// GetLightclientAuthorityMock returns a new lightclient authority keeper mock
func GetLightclientAuthorityMock(t testing.TB, keeper *keeper.Keeper) *lightclientmocks.LightclientAuthorityKeeper {
	cok, ok := keeper.GetAuthorityKeeper().(*lightclientmocks.LightclientAuthorityKeeper)
	require.True(t, ok)
	return cok
}
