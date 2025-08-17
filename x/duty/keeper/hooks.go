package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/duty/types"
)

// Hooks wrapper struct for duty keeper
type Hooks struct {
	k Keeper
}

var _ types.DutyHooks = Hooks{}

// Hooks returns the duty keeper hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// AfterDutyCreated is called after a duty is created
func (h Hooks) AfterDutyCreated(ctx sdk.Context, duty types.Duty) error {
	// Add any logic that should run after a duty is created
	return nil
}

// AfterDutyUpdated is called after a duty is updated
func (h Hooks) AfterDutyUpdated(ctx sdk.Context, duty types.Duty) error {
	// Add any logic that should run after a duty is updated
	return nil
}

// AfterDutyDeleted is called after a duty is deleted
func (h Hooks) AfterDutyDeleted(ctx sdk.Context, duty types.Duty) error {
	// Add any logic that should run after a duty is deleted
	return nil
}

