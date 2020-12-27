package teststore

import (
	"github.com/zhoreeq/ipfd/internal/app/models"
	"github.com/zhoreeq/ipfd/internal/app/store"
)

func New() *Store {
	return &Store{}
}

type Store struct {
	addressRepository *AddressRepository
	postRepository    *PostRepository
	commentRepository *CommentRepository
	voteRepository    *VoteRepository
}

func (s *Store) Address() store.AddressRepository {
	if s.addressRepository != nil {
		return s.addressRepository
	}

	s.addressRepository = &AddressRepository{
		store: s,
		addresses: make(map[int]*models.Address),
	}

	return s.addressRepository
}

func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}

	s.postRepository = &PostRepository{
		store: s,
		posts: make(map[int]*models.Post),
	}

	return s.postRepository
}

func (s *Store) Comment() store.CommentRepository {
	if s.commentRepository != nil {
		return s.commentRepository
	}

	s.commentRepository = &CommentRepository{
		store: s,
		comments: make(map[int]*models.Comment),
	}

	return s.commentRepository
}

func (s *Store) Vote() store.VoteRepository {
	if s.voteRepository != nil {
		return s.voteRepository
	}

	s.voteRepository = &VoteRepository{
		store: s,
	}

	return s.voteRepository
}
