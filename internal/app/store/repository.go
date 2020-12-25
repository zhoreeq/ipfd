package store

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type AddressRepository interface {
	GetOrCreate(m *models.Address) error
}

type PostRepository interface {
	Create(m *models.Post) error
	GetBoardIndex(page, limit int, published bool) ([]*models.Post, error)
	GetById(id string) (*models.Post, error)
}

type CommentRepository interface {
	Create(m *models.Comment) error
	GetByPostId(postId string) ([]*models.Comment, error)
	GetAll(page, limit int) ([]*models.Comment, error)
}

type VoteRepository interface {
	Upvote(m *models.Vote) error
	Downvote(m *models.Vote) error
}
