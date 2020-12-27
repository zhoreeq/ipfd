package teststore

import (
	"strconv"

	"github.com/zhoreeq/ipfd/internal/app/models"
)

type CommentRepository struct {
	store *Store
	comments map[int]*models.Comment
}

func (r *CommentRepository) Create(m *models.Comment) error {
	m.Id = len(r.comments) + 1
	r.comments[m.Id] = m
	return nil
}

func (r *CommentRepository) GetByPostId(postId string) ([]*models.Comment, error) {
	var postIdInt int
	var err error
	if postIdInt, err = strconv.Atoi(postId); err != nil {
		return nil, err
	}

	var comments []*models.Comment
	for i, c := range r.comments {
		if c.PostId == postIdInt {
			c.Id = i
			comments = append(comments, c)
		}
	}
	return comments, nil

}

func (r *CommentRepository) GetAll(page, limit int) ([]*models.Comment, error) {
	var allComments []*models.Comment

	for i, c := range r.comments {
		c.Id = i
		allComments = append(allComments, c)
	}
	return allComments, nil

}
