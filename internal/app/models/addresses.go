package models

import (
	"net"
)

type Address struct {
	Id        int
	Address   *net.IP
	CreatedAt string
}
