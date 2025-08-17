package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// We don't need to write anything on-chain for updates; we emit events
// so off-chain indexers / agents can react.
type DutyHooks struct{ k Keeper }

var _ stakingtypes.StakingHooks = DutyHooks{}

func (h DutyHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	// Get validator info for more detailed event
	validator, found := h.k.stakingKeeper.GetValidator(ctx, valAddr)
	if found {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("duty_validator_bonded",
				sdk.NewAttribute("cons_addr", consAddr.String()),
				sdk.NewAttribute("val_addr", valAddr.String()),
				sdk.NewAttribute("voting_power", validator.GetTokens().String()),
				sdk.NewAttribute("moniker", validator.GetMoniker()),
			),
		)
	} else {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("duty_validator_bonded",
				sdk.NewAttribute("cons_addr", consAddr.String()),
				sdk.NewAttribute("val_addr", valAddr.String()),
			),
		)
	}
}

func (h DutyHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_validator_removed",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("val_addr", valAddr.String()),
		),
	)
}

func (h DutyHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_validator_unbonding",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("val_addr", valAddr.String()),
		),
	)
}

// Implement other hooks as no-ops for brevity
func (h DutyHooks) AfterValidatorCreated(sdk.Context, sdk.ValAddress)                          {}
func (h DutyHooks) BeforeValidatorModified(sdk.Context, sdk.ValAddress)                        {}
func (h DutyHooks) BeforeDelegationCreated(sdk.Context, sdk.ValAddress, sdk.AccAddress)        {}
func (h DutyHooks) AfterDelegationModified(sdk.Context, sdk.ValAddress, sdk.AccAddress)        {}
func (h DutyHooks) BeforeDelegationSharesModified(sdk.Context, sdk.ValAddress, sdk.AccAddress) {}
func (h DutyHooks) AfterUnbondingInitiated(sdk.Context, uint64)                                {}
