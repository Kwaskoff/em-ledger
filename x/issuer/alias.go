package issuer

import (
	"github.com/e-money/em-ledger/x/issuer/keeper"
	"github.com/e-money/em-ledger/x/issuer/types"
)

const (
	StoreKey   = types.StoreKey
	ModuleName = types.ModuleName
)

var (
	ModuleCdc = types.ModuleCdc
	NewKeeper = keeper.NewKeeper
	NewIssuer = types.NewIssuer
)

type (
	Keeper = keeper.Keeper
	Issuer = types.Issuer
)
