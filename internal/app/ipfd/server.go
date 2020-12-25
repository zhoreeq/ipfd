package ipfd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"

	"github.com/zhoreeq/ipfd/internal/app/ipfs"
	"github.com/zhoreeq/ipfd/internal/app/store"
)

func New(config *Config, log *log.Logger, dbStore store.Store, ipfsShells []ipfs.Shell) *Server {
	server := &Server{}
	server.config = config
	server.log = log
	server.router = mux.NewRouter()
	server.handler = alice.New(nosurf.NewPure).Then(server.router)
	server.store = dbStore
	server.ipfs = ipfsShells

	server.configureTemplates()
	server.configureRouter()

	return server
}

type Server struct {
	config   *Config
	log      *log.Logger
	router   *mux.Router
	handler  http.Handler
	store    store.Store
	template *template.Template
	ipfs     []ipfs.Shell
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/", s.BoardHandler)
	s.router.HandleFunc("/create_post", s.CreatePostHandler).Methods("POST")
	s.router.HandleFunc("/res/{id:[0-9]+}.html", s.PostViewHandler)
	s.router.HandleFunc("/board.rss", s.PostFeed)
	s.router.HandleFunc("/comments.rss", s.CommentFeed)

	if s.config.EnableComments {
		s.router.HandleFunc("/create_comment", s.CreateCommentHandler).Methods("POST")
	}
	if s.config.EnableVotes {
		s.router.HandleFunc("/upvote/{id:[0-9]+}", s.UpvoteHandler)
		s.router.HandleFunc("/downvote/{id:[0-9]+}", s.DownvoteHandler)
	}
	if s.config.ServeStatic {
		s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.config.StaticPath))))
	}
}

func (s *Server) configureTemplates() {
	templatesGlob := fmt.Sprintf("%s/*.html", s.config.TemplatesPath)
	funcs := template.FuncMap{
		"mediaTag":    s.mediaTag,
		"staticURL":   s.getStaticUrl,
		"getSiteName": s.getSiteName,
	}
	s.template = template.Must(template.New("ipfd").Funcs(funcs).ParseGlob(templatesGlob))
}

func (s *Server) getIpfsUrl(cid string) string {
	return fmt.Sprintf("%s%s", s.config.IpfsGateway, cid)
}

func (s *Server) getStaticUrl(url string) string {
	return fmt.Sprintf("%s%s", s.config.StaticURL, url)
}

func (s *Server) getSiteName() string {
	return s.config.SiteName
}

func (s *Server) mediaTag(cid, contentType string) template.HTML {
	var content string
	if strings.HasPrefix(contentType, "image/") {
		content = fmt.Sprintf("<img src=\"%s%s\">", s.config.IpfsGateway, cid)
	} else if strings.HasPrefix(contentType, "video/") {
		content = fmt.Sprintf("<video src=\"%s%s\" controls loop></video>", s.config.IpfsGateway, cid)
	} else if strings.HasPrefix(contentType, "audio/") {
		content = fmt.Sprintf("<audio src=\"%s%s\" controls loop></audio>", s.config.IpfsGateway, cid)
	} else {
		content = fmt.Sprintf("<a href=\"%s%s\" target=\"_blank\">download</a>", s.config.IpfsGateway, cid)
	}
	return template.HTML(content)
}
