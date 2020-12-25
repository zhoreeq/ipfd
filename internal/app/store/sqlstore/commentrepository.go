package sqlstore

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type CommentRepository struct {
	store *Store
}

func (r *CommentRepository) Create(m *models.Comment) error {
	return r.store.db.QueryRow(
		"INSERT INTO comments (address_id, post_id, text) VALUES ($1, $2, $3) RETURNING id",
		m.AddressId,
		m.PostId,
		m.Text,
	).Scan(&m.Id)
}

func (r *CommentRepository) GetByPostId(postId string) ([]*models.Comment, error) {
	rows, err := r.store.db.Query("SELECT id, text FROM comments WHERE post_id = $1 ORDER BY id asc", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.Id, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepository) GetAll(page, limit int) ([]*models.Comment, error) {
	offset := 0
	if page > 0 {
		offset = page * limit
	}

	rows, err := r.store.db.Query(`SELECT id, text, post_id, created_at 
			FROM comments ORDER BY id desc LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}
	for rows.Next() {
		comment := &models.Comment{}
		err := rows.Scan(&comment.Id, &comment.Text, &comment.PostId, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
