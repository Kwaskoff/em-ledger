// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const ClientOrderIDMaxLength = 32

var (
	_ sdk.Msg = MsgAddLimitOrder{}
	_ sdk.Msg = MsgAddMarketOrder{}
	_ sdk.Msg = MsgCancelOrder{}
	_ sdk.Msg = MsgCancelReplaceOrder{}
)

type (
	MsgAddLimitOrder struct {
		Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
		TimeInForce   string         `json:"time_in_force" yaml:"time_in_force"`
		Source        sdk.Coin       `json:"source" yaml:"source"`
		Destination   sdk.Coin       `json:"destination" yaml:"destination"`
		ClientOrderId string         `json:"client_order_id" yaml:"client_order_id"`
	}

	MsgAddMarketOrder struct {
		Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
		TimeInForce   string         `json:"time_in_force" yaml:"time_in_force"`
		Source        string         `json:"source" yaml:"source"`
		Destination   sdk.Coin       `json:"destination" yaml:"destination"`
		ClientOrderId string         `json:"client_order_id" yaml:"client_order_id"`
		MaxSlippage   sdk.Dec        `json:"maximum_slippage" yaml:"maximum_slippage"`
	}

	MsgCancelOrder struct {
		Owner         sdk.AccAddress `json:"owner" yaml:"owner"`
		ClientOrderId string         `json:"client_order_id" yaml:"client_order_id"`
	}

	MsgCancelReplaceOrder struct {
		Owner             sdk.AccAddress `json:"owner" yaml:"owner"`
		Source            sdk.Coin       `json:"source" yaml:"source"`
		Destination       sdk.Coin       `json:"destination" yaml:"destination"`
		OrigClientOrderId string         `json:"original_client_order_id" yaml:"original_client_order_id"`
		NewClientOrderId  string         `json:"new_client_order_id" yaml:"new_client_order_id"`
		MaxSlippage       sdk.Dec        `json:"maximum_slippage" yaml:"maximum_slippage"`
	}
)

func (m MsgAddMarketOrder) Route() string {
	return RouterKey
}

func (m MsgAddMarketOrder) Type() string {
	return "add_market_order"
}

func (m MsgAddMarketOrder) ValidateBasic() error {
	if m.MaxSlippage.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(ErrInvalidSlippage, "Cannot be negative")
	}

	if m.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}

	if !m.Destination.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "destination amount is invalid: %v", m.Destination.String())
	}

	err := sdk.ValidateDenom(m.Source)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "source denomination is invalid: %v", m.Source)
	}

	if m.Source == m.Destination.Denom {
		return sdkerrors.Wrapf(ErrInvalidInstrument, "'%v/%v' is not a valid instrument", m.Source, m.Destination.Denom)
	}

	return validateClientOrderID(m.ClientOrderId)

}

func (m MsgAddMarketOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAddMarketOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

func (m MsgCancelReplaceOrder) Route() string {
	return RouterKey
}

func (m MsgCancelReplaceOrder) Type() string {
	return "cancel_replace_order"
}

func (m MsgCancelReplaceOrder) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}

	if !m.Destination.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "destination amount is invalid: %v", m.Destination.String())
	}

	if !m.Source.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "source amount is invalid: %v", m.Source.String())
	}

	if m.Source.Denom == m.Destination.Denom {
		return sdkerrors.Wrapf(ErrInvalidInstrument, "'%v/%v' is not a valid instrument", m.Source.Denom, m.Destination.Denom)
	}

	err := validateClientOrderID(m.OrigClientOrderId)
	if err != nil {
		return err
	}

	return validateClientOrderID(m.NewClientOrderId)
}

func (m MsgCancelReplaceOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCancelReplaceOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

func (m MsgCancelOrder) Route() string {
	return RouterKey
}

func (m MsgCancelOrder) Type() string {
	return "cancel_order"
}

func (m MsgCancelOrder) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}

	return validateClientOrderID(m.ClientOrderId)
}

func (m MsgCancelOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCancelOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

func (m MsgAddLimitOrder) Route() string {
	return RouterKey
}

func (m MsgAddLimitOrder) Type() string {
	return "add_limit_order"
}

func (m MsgAddLimitOrder) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}

	if !m.Destination.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "destination amount is invalid: %v", m.Destination.String())
	}

	if !m.Source.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "source amount is invalid: %v", m.Source.String())
	}

	if m.Source.Denom == m.Destination.Denom {
		return sdkerrors.Wrapf(ErrInvalidInstrument, "'%v/%v' is not a valid instrument", m.Source.Denom, m.Destination.Denom)
	}

	return validateClientOrderID(m.ClientOrderId)
}

func (m MsgAddLimitOrder) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAddLimitOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

func validateClientOrderID(id string) error {
	if len(id) > ClientOrderIDMaxLength {
		return sdkerrors.Wrap(ErrInvalidClientOrderId, id)
	}

	return nil
}
