package teststore

import (
	"strconv"

	"github.com/zhoreeq/ipfd/internal/app/models"
)

type PostRepository struct {
	store *Store
	posts map[int]*models.Post
}

func (r *PostRepository) Create(m *models.Post) error {
	m.Id = len(r.posts) + 1
	r.posts[m.Id] = m
	return nil
}

func (r *PostRepository) GetBoardIndex(page, limit int, published bool) ([]*models.Post, error) {
	var allPosts []*models.Post
	for i, p := range r.posts {
		p.Id = i
		allPosts = append(allPosts, p)
	}
	return allPosts, nil
}

func (r *PostRepository) GetById(id string) (*models.Post, error) {
	var idInt int
	var err error
	if idInt, err = strconv.Atoi(id); err != nil {
		return nil, err
	}
	p := r.posts[idInt]
	return p, nil
}
