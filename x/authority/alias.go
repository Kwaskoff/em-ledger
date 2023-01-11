package authority

import (
	"github.com/e-money/em-ledger/x/authority/client/cli"
	"github.com/e-money/em-ledger/x/authority/keeper"
	"github.com/e-money/em-ledger/x/authority/types"
)

const (
	ModuleName     = types.ModuleName
	StoreKey       = types.StoreKey
	QuerierRoute   = types.QuerierRoute
	QueryGasPrices = types.QueryGasPrices
)

type (
	Keeper = keeper.Keeper
)

var (
	ModuleCdc       = types.ModuleCdc
	NewKeeper       = keeper.NewKeeper
	BeginBlocker    = keeper.BeginBlocker
	GetGasPricesCmd = cli.GetGasPricesCmd
	GetQueryCmd     = cli.GetQueryCmd
	GetTxCmd        = cli.GetTxCmd
)
