package types

import (
	"fmt"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId: PortID,
		// this line is used by starport scaffolding # genesis/types/default
		IbcTokenList: []*IbcToken{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}

	// this line is used by starport scaffolding # genesis/types/validate
	// Check for duplicated index in ibcToken
	ibcTokenIndexMap := make(map[string]struct{})

	for _, elem := range gs.IbcTokenList {
		index := string(IbcTokenKey(elem.Index))
		if _, ok := ibcTokenIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for ibcToken")
		}
		ibcTokenIndexMap[index] = struct{}{}
	}

	return nil
}
