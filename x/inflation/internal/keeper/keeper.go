package keeper

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"

	"emoney/x/inflation/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper of the mint store
type Keeper struct {
	cdc              *codec.Codec
	storeKey         sdk.StoreKey
	paramSpace       params.Subspace
	supplyKeeper     types.SupplyKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	supplyKeeper types.SupplyKeeper, feeCollectorName string) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// get the minter
func (k Keeper) GetState(ctx sdk.Context) (is types.InflationState) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored inflation state should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &is)
	return
}

// set the minter
func (k Keeper) SetState(ctx sdk.Context, is types.InflationState) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(is)
	store.Set(types.MinterKey, b)
}

//______________________________________________________________________

// GetParams returns the total set of minting parameters.
//func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
//	k.paramSpace.Get(ctx, types.KeyParams, &params)
//	return params
//}
//
//// SetParams sets the total set of minting parameters.
//func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
//	k.paramSpace.Set(ctx, types.KeyParams, params)
//}

//______________________________________________________________________

func (k Keeper) TotalTokenSupply(ctx sdk.Context) sdk.Coins {
	return k.supplyKeeper.GetSupply(ctx).Total
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) sdk.Error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

func (k Keeper) AddMintedCoins(ctx sdk.Context, fees sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
