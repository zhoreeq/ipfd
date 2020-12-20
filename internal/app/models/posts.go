package models

import (
	"time"
)

type Post struct {
	Id            int
	AddressId     int
	Title         string
	CID           string
	FileSize      int64
	ContentType   string
	CreatedAt     time.Time
	Published     bool
	CommentsCount int
	UpCount       int
	DownCount     int
}
