package models

import (
	"time"
)

type Comment struct {
	Id        int
	PostId    int
	AddressId int
	Text      string
	CreatedAt time.Time
}
