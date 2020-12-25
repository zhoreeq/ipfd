package store

type Store interface {
	Address() AddressRepository
	Post() PostRepository
	Comment() CommentRepository
	Vote() VoteRepository
}
