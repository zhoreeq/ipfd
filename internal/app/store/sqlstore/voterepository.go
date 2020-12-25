package sqlstore

import (
	// "database/sql"

	"github.com/zhoreeq/ipfd/internal/app/models"
)

type VoteRepository struct {
	store *Store
}

func (r *VoteRepository) Upvote(m *models.Vote) error {
	return r.store.db.QueryRow(
		"INSERT INTO upvotes (address_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id",
		m.AddressId,
		m.PostId,
	).Scan(&m.Id)
}

func (r *VoteRepository) Downvote(m *models.Vote) error {
	return r.store.db.QueryRow(
		"INSERT INTO downvotes (address_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id",
		m.AddressId,
		m.PostId,
	).Scan(&m.Id)
}
