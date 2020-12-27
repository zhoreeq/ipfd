package ipfs

import (
	"io"

	ipfsApi "github.com/ipfs/go-ipfs-api"
)

type FakeShell struct {}

func (s *FakeShell) Add(r io.Reader, options ...ipfsApi.AddOpts) (string, error){
	return "FAKE", nil
}

func (s *FakeShell) Pin(path string) error {
	return nil
}

