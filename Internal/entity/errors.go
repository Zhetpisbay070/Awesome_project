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
	PozdnoNahui               = errors.New("uzhe ne otmenit")
)
