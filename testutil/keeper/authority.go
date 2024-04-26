package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/keeper"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

var (
	AuthorityGovAddress = sample.Bech32AccAddress()
)

func initAuthorityKeeper(
	cdc codec.Codec,
	db *tmdb.MemDB,
	ss store.CommitMultiStore,
) keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	ss.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	ss.MountStoreWithDB(memKey, storetypes.StoreTypeMemory, db)

	return keeper.NewKeeper(
		cdc,
		storeKey,
		memKey,
		AuthorityGovAddress,
	)
}

// AuthorityKeeper instantiates an authority keeper for testing purposes
func AuthorityKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	// Initialize local store
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	cdc := NewCodec()

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

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		AuthorityGovAddress,
	)

	return &k, ctx
}

// MockIsAuthorized mocks the IsAuthorized method of an authority keeper mock
func MockIsAuthorized(m *mock.Mock, address string, policyType types.PolicyType, isAuthorized bool) {
	m.On("IsAuthorized", mock.Anything, address, policyType).Return(isAuthorized).Once()
}

func SetAdminPolices(ctx sdk.Context, ak *keeper.Keeper) string {
	admin := sample.AccAddress()
	ak.SetPolicies(ctx, types.Policies{Items: []*types.Policy{
		{
			Address:    admin,
			PolicyType: types.PolicyType_groupAdmin,
		},
	}})
	return admin
}
