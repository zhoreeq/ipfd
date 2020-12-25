package sqlstore

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type AddressRepository struct {
	store *Store
}

func (r *AddressRepository) GetOrCreate(m *models.Address) error {
	return r.store.db.QueryRow(
		"INSERT INTO addresses (address) VALUES ($1) ON CONFLICT (address) DO UPDATE SET created_at = Now() RETURNING id",
		m.Address,
	).Scan(&m.Id)
}
