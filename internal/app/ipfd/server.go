package ipfd

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/justinas/alice"
	_ "github.com/jackc/pgx/v4/stdlib"

	ipfsApi "github.com/ipfs/go-ipfs-api"

	"github.com/zhoreeq/ipfd/internal/app/ipfs"
	"github.com/zhoreeq/ipfd/internal/app/store"
)

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func newIpfsShell(ipfsAPI string) (*ipfsApi.Shell, error) {
	shell := ipfsApi.NewShell(ipfsAPI)
	if _, _, err := shell.Version(); err != nil {
		return nil, err
	}

	return shell, nil
}

type Server struct {
	config   *Config
	log      *log.Logger
	router   *mux.Router
	store    *store.Store
	template *template.Template
	ipfs     ipfs.Shell
}

func New(config *Config, log *log.Logger) *Server {
	server := &Server{}
	server.config = config
	server.log = log
	server.router = mux.NewRouter()

	return server
}

func Start(config *Config, log *log.Logger) error {
	server := New(config, log)

	db, err := newDB(server.config.databaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	server.store = store.New(db)

	shell, err := newIpfsShell(server.config.ipfsAPI)
	if err != nil {
		return err
	}
	server.ipfs = shell

	server.configureTemplates()
	server.configureRouter()
	chain := alice.New(nosurf.NewPure).Then(server.router)
	http.Handle("/", chain)
	server.log.Println("Started on", server.config.bindAddress)

	return http.ListenAndServe(server.config.bindAddress, nil)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/", s.BoardHandler)
	s.router.HandleFunc("/create_post", s.CreatePostHandler).Methods("POST")
	s.router.HandleFunc("/res/{id:[0-9]+}.html", s.PostViewHandler)
	s.router.HandleFunc("/board.rss", s.PostFeed)
	s.router.HandleFunc("/comments.rss", s.CommentFeed)

	if s.config.enableComments {
		s.router.HandleFunc("/create_comment", s.CreateCommentHandler).Methods("POST")
	}
	if s.config.enableVotes {
		s.router.HandleFunc("/upvote/{id:[0-9]+}", s.UpvoteHandler)
		s.router.HandleFunc("/downvote/{id:[0-9]+}", s.DownvoteHandler)
	}
	if s.config.serveStatic {
		s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.config.staticPath))))
	}
}

func (s *Server) configureTemplates() {
	templatesGlob := fmt.Sprintf("%s/*.html", s.config.templatesPath)
	funcs := template.FuncMap{
		"mediaTag":    s.mediaTag,
		"staticURL":   s.getStaticUrl,
		"getSiteName": s.getSiteName,
	}
	s.template = template.Must(template.New("ipfd").Funcs(funcs).ParseGlob(templatesGlob))
}

func (s *Server) getIpfsUrl(cid string) string {
	return fmt.Sprintf("%s%s", s.config.ipfsGateway, cid)
}

func (s *Server) getStaticUrl(url string) string {
	return fmt.Sprintf("%s%s", s.config.staticURL, url)
}

func (s *Server) getSiteName() string {
	return s.config.siteName
}

func (s *Server) mediaTag(cid, contentType string) template.HTML {
	var content string
	if strings.HasPrefix(contentType, "image/") {
		content = fmt.Sprintf("<img src=\"%s%s\">", s.config.ipfsGateway, cid)
	} else if strings.HasPrefix(contentType, "video/") {
		content = fmt.Sprintf("<video src=\"%s%s\" controls loop></video>", s.config.ipfsGateway, cid)
	} else if strings.HasPrefix(contentType, "audio/") {
		content = fmt.Sprintf("<audio src=\"%s%s\" controls loop></audio>", s.config.ipfsGateway, cid)
	} else {
		content = fmt.Sprintf("<a href=\"%s%s\" target=\"_blank\">download</a>", s.config.ipfsGateway, cid)
	}
	return template.HTML(content)
}
