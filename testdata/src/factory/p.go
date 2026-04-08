package factory

import "domain"

func NewAddress(val string) (domain.Address, error) {
	return domain.Address(val), nil
}
