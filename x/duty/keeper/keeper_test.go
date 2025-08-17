package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/TheArticulation/Duty/x/duty/types"
)

func setupTestKeeper(t *testing.T) (Keeper, sdk.Context) {
	// Create a test database
	db := tmdb.NewMemDB()

	// Create a test store
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(types.StoreKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	// Create a test context
	ctx := sdk.NewContext(ms, tmproto.Header{}, false, log.NewNopLogger())

	// Create a test codec
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)

	// Create a test param space
	paramSpace := paramtypes.NewSubspace(cdc, types.ModuleName, "duty", "duty")

	// Create a test keeper
	keeper := NewKeeper(
		cdc,
		store.NewKVStoreService(ms.GetKVStore(types.StoreKey)),
		paramSpace,
		nil, // staking keeper (nil for test)
		log.NewNopLogger(),
		nil, // params service (nil for test)
	)

	return keeper, ctx
}

func TestKeeper_SetAndGetDutyMetadata(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)

	// Create a test consensus address
	consAddr := sdk.ConsAddress([]byte("test-validator"))

	// Create test metadata
	metadata := types.DutyMetadata{
		CheckpointPubKey:     "0x1234567890abcdef",
		CheckpointStorageURI: "s3://bucket/prefix/",
	}

	// Set metadata
	keeper.SetDutyMetadata(ctx, consAddr, metadata)

	// Get metadata
	retrievedMetadata, found := keeper.GetDutyMetadata(ctx, consAddr)

	// Assertions
	assert.True(t, found)
	assert.Equal(t, metadata.CheckpointPubKey, retrievedMetadata.CheckpointPubKey)
	assert.Equal(t, metadata.CheckpointStorageURI, retrievedMetadata.CheckpointStorageURI)
}

func TestKeeper_GetParams(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)

	// Get default params
	params := keeper.GetParams(ctx)

	// Assert default values
	assert.Equal(t, uint32(2), params.QuorumNumerator)
	assert.Equal(t, uint32(3), params.QuorumDenominator)
}

func TestKeeper_SetParams(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)

	// Create custom params
	customParams := types.Params{
		QuorumNumerator:   3,
		QuorumDenominator: 4,
	}

	// Set params
	keeper.SetParams(ctx, customParams)

	// Get params
	retrievedParams := keeper.GetParams(ctx)

	// Assertions
	assert.Equal(t, customParams.QuorumNumerator, retrievedParams.QuorumNumerator)
	assert.Equal(t, customParams.QuorumDenominator, retrievedParams.QuorumDenominator)
}

func TestKeeper_GetDutySet(t *testing.T) {
	keeper, ctx := setupTestKeeper(t)

	// Get duty set (without staking keeper, this will return empty list)
	validators, params := keeper.GetDutySet(ctx)

	// Assertions
	assert.Empty(t, validators) // No validators without staking keeper
	assert.Equal(t, uint32(2), params.QuorumNumerator)
	assert.Equal(t, uint32(3), params.QuorumDenominator)
}

func TestKeeper_ParamsValidation(t *testing.T) {
	// Test valid params
	validParams := types.Params{
		QuorumNumerator:   2,
		QuorumDenominator: 3,
	}
	assert.NoError(t, validParams.Validate())

	// Test invalid params - zero numerator
	invalidParams1 := types.Params{
		QuorumNumerator:   0,
		QuorumDenominator: 3,
	}
	assert.Error(t, invalidParams1.Validate())

	// Test invalid params - zero denominator
	invalidParams2 := types.Params{
		QuorumNumerator:   2,
		QuorumDenominator: 0,
	}
	assert.Error(t, invalidParams2.Validate())

	// Test invalid params - numerator > denominator
	invalidParams3 := types.Params{
		QuorumNumerator:   4,
		QuorumDenominator: 3,
	}
	assert.Error(t, invalidParams3.Validate())
}
