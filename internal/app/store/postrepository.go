package store

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
)

type PostRepository struct {
	store *Store
}

func (r *PostRepository) Create(m *models.Post) error {
	return r.store.db.QueryRow(
		"INSERT INTO posts (address_id, title, cid, file_size, content_type, published) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		m.AddressId,
		m.Title,
		m.CID,
		m.FileSize,
		m.ContentType,
		m.Published,
	).Scan(&m.Id)
}

func (r *PostRepository) GetBoardIndex(page, limit int, published bool) ([]*models.Post, error) {
	offset := 0
	if page > 0 {
		offset = page * limit
	}

	rows, err := r.store.db.Query(`SELECT id, title, cid, content_type, created_at,
			(select count(id) from comments where comments.post_id = posts.id) as comments_count,
			(select count(id) from upvotes where upvotes.post_id = posts.id) as upvotes_count,
			(select count(id) from downvotes where downvotes.post_id = posts.id) as downvotes_count 
			FROM posts WHERE published = $1 ORDER BY id desc LIMIT $2 OFFSET $3`, published, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*models.Post{}

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(&post.Id, &post.Title, &post.CID, &post.ContentType, &post.CreatedAt, &post.CommentsCount, &post.UpCount, &post.DownCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}

	return posts, nil
}

func (r *PostRepository) GetById(id string) (*models.Post, error) {
	post := &models.Post{}
	err := r.store.db.QueryRow(`select id, title, cid, content_type,
		(select count(id) from upvotes where upvotes.post_id = posts.id) as upvotes_count,
		(select count(id) from downvotes where downvotes.post_id = posts.id) as downvotes_count 
		FROM posts where id = $1 and published = true`, id).Scan(&post.Id, &post.Title, &post.CID, &post.ContentType, &post.UpCount, &post.DownCount)
	if err != nil {
		return nil, err
	}

	return post, nil
}
