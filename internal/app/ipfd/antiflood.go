package ipfd

import (
	"net/http"
)

func (s *Server) isPostingAllowed(r *http.Request) bool {
	return true
}
