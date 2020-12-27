package teststore

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type AddressRepository struct {
	store *Store
	addresses map[int]*models.Address
}

func (r *AddressRepository) GetOrCreate(m *models.Address) error {
	return nil
}
