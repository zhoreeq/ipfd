package store

import (
	"database/sql"
)

type Store struct {
	db                *sql.DB
	addressRepository *AddressRepository
	postRepository    *PostRepository
	commentRepository *CommentRepository
	voteRepository    *VoteRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Address() *AddressRepository {
	if s.addressRepository != nil {
		return s.addressRepository
	}

	s.addressRepository = &AddressRepository{
		store: s,
	}

	return s.addressRepository
}

func (s *Store) Post() *PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
	}

	return s.postRepository
}

func (s *Store) Comment() *CommentRepository {
	if s.commentRepository != nil {
		return s.commentRepository
	}

	s.commentRepository = &CommentRepository{
		store: s,
	}

	return s.commentRepository
}

func (s *Store) Vote() *VoteRepository {
	if s.voteRepository != nil {
		return s.voteRepository
	}

	s.voteRepository = &VoteRepository{
		store: s,
	}

	return s.voteRepository
}
