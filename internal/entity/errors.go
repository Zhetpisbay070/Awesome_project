package entity

import "errors"

var (
	ProductDoesNotExistError  = errors.New("product does not exist")
	OrderCannotBeCancelled    = errors.New("you cannot cancel your order at this stage")
	InvalidStatus             = errors.New("no status")
	ImpossibleToCheckProducts = errors.New("impossible to check products")
	DelTypeUnavailable        = errors.New("delivery type unavailable")
	OrderNotFound             = errors.New("order not found")
	InvalidTransition         = errors.New("invalid transition")
	Pozdno                    = errors.New("uzhe ne otmenit")
	OrderCannotBeEdited       = errors.New("order cannot be edited in this current status")
	AddressCannotBeEdited     = errors.New("address cannot be edited")
	OrderCannotBeChanged      = errors.New("order cannot be changed")
	ErrOrderNotFound          = errors.New("order not found")
)
