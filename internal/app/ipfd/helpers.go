package ipfd

import (
	"net"
	"net/http"
	"strconv"

	"github.com/zhoreeq/ipfd/internal/app/models"
)

func getRealIP(r *http.Request) (*net.IP, error) {
	if realIP, ok := r.Header["X-Real-Ip"]; ok && len(r.Header["X-Real-Ip"]) == 1 {
		ip := net.ParseIP(realIP[0])
		return &ip, nil
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(host)

	return &ip, nil
}

func getPage(r *http.Request) int {
	params := r.URL.Query()
	if _, ok := params["page"]; ok {
		if page, err := strconv.Atoi(params["page"][0]); err == nil {
			return page
		}
	}
	return 0
}

func (s *Server) GetAddress(r *http.Request) (*models.Address, error) {
	addr := &models.Address{}
	if remoteIP, err := getRealIP(r); err != nil {
		return nil, err
	} else {
		addr.Address = remoteIP
	}

	if err := s.store.Address().GetOrCreate(addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (s *Server) RedirectBack(w http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Referer"]; ok && len(r.Header["Referer"]) == 1 {
		w.Header().Set("Location", r.Header["Referer"][0])
		w.WriteHeader(303)
	} else {
		http.Redirect(w, r, "./", 303)
	}
}

// Log an error and display a generic message
func (s *Server) InternalError(w *http.ResponseWriter, err error) {
	s.log.Println(err.Error())
	http.Error(*w, "Internal server error", 500)
}
