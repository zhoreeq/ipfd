package ipfd

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zhoreeq/ipfd/internal/app/ipfs"
	"github.com/zhoreeq/ipfd/internal/app/store"
	"github.com/zhoreeq/ipfd/internal/app/store/teststore"
)

func newTestServer(store store.Store) *Server {
	logger := log.New(os.Stdout, "", log.Flags())

	conf := &Config{
		SiteURL:             "http://127.0.0.1:8000",
		SitePath:            "/",
		SiteName:            "test ipfd",
		BindAddress:         "[::]:8001",
		DatabaseURL:         "",
		TemplatesPath:       "../../../templates",
		StaticPath:          "./static",
		StaticURL:           "/static",
		ServeStatic:         false,
		IpfsAPI:             []string{"/ip4/127.0.0.1/tcp/5001"},
		IpfsGateway:         "http://127.0.0.1:8080/ipfs",
		IpfsPin:             true,
		MaxFileSize:         6000000,
		AllowedContentTypes: []string{"image/gif", "image/png", "image/jpeg", "video/mp4", "video/webm", "audio/ogg", "audio/mpeg"},
		Premoderation:       false,
		EnableComments:      true,
		EnableVotes:         true,
	}

	var ipfsShells []ipfs.Shell
	ipfsShells = append(ipfsShells, &ipfs.FakeShell{})

	s := New(conf, logger, store, ipfsShells)
	return s
}

func TestIpfdServer_Index(t *testing.T) {
	store := teststore.New()
	s := newTestServer(store)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	s.ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK, "should be equal")
}
