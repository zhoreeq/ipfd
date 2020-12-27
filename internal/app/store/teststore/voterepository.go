package teststore

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type VoteRepository struct {
	store *Store
}

func (r *VoteRepository) Upvote(m *models.Vote) error {
	return nil
}

func (r *VoteRepository) Downvote(m *models.Vote) error {
	return nil
}
