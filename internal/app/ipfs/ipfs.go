package ipfs

import (
	"io"

	ipfsApi "github.com/ipfs/go-ipfs-api"
)

type Shell interface {
	Add(r io.Reader, options ...ipfsApi.AddOpts) (string, error)
}
