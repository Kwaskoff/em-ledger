package types

import (
	"fmt"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewIssuer(address sdk.AccAddress, denoms ...string) Issuer {
	sort.Strings(denoms)

	return Issuer{
		Address: address.String(),
		Denoms:  denoms,
	}
}

func (i Issuer) IsValid() bool {
	if len(i.Address) == 0 {
		return false
	}

	if len(i.Denoms) == 0 {
		return false
	}

	return true
}

func (i Issuers) String() string {
	var sb strings.Builder

	for _, issuer := range i.Issuers {
		sb.WriteString(fmt.Sprintf("%v : %v\n", issuer.Address, issuer.Denoms))
	}

	return sb.String()
}
