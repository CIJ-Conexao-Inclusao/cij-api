package service

import (
	"cij_api/src/model"
	"cij_api/src/repo"
	"cij_api/src/utils"
)

type AddressService interface {
	GetAddressById(id int) (model.Address, utils.Error)
}

type addressService struct {
	addressRepo repo.AddressRepo
}

func NewAddressService(addressRepo repo.AddressRepo) AddressService {
	return &addressService{
		addressRepo: addressRepo,
	}
}

func (n *addressService) GetAddressById(id int) (model.Address, utils.Error) {
	address, err := n.addressRepo.GetAddressById(id)
	if err.Code != "" {
		return address, err
	}

	return address, utils.Error{}
}
