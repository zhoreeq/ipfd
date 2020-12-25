package ipfd

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	ipfs "github.com/ipfs/go-ipfs-api"
	"github.com/justinas/nosurf"

	"github.com/zhoreeq/ipfd/internal/app/models"
)

// Add pin CID at a specified s.ipfs shell
func (s *Server) pinAtShell(shellId int, cid string) {
	err := s.ipfs[shellId].Pin(cid)
	if err != nil {
		s.log.Println("Pin error: ", err)
	}
}

func (s *Server) BoardHandler(w http.ResponseWriter, r *http.Request) {
	page := getPage(r)

	posts, err := s.store.Post().GetBoardIndex(page, 10, true)
	if err != nil {
		s.InternalError(&w, err)
		return
	}
	postingAllowed := s.isPostingAllowed(r)

	var prevPage, nextPage int
	if len(posts) == 10 {
		nextPage = page + 1
	}
	if page > 0 {
		prevPage = page - 1
	}

	data := struct {
		NextPage       int
		PrevPage       int
		Posts          []*models.Post
		PostingAllowed bool
		EnableVotes    bool
		Token          string
	}{
		NextPage:       nextPage,
		PrevPage:       prevPage,
		Posts:          posts,
		PostingAllowed: postingAllowed,
		EnableVotes:    s.config.EnableVotes,
		Token:          nosurf.Token(r),
	}

	err = s.template.ExecuteTemplate(w, "board.html", data)
	if err != nil {
		s.InternalError(&w, err)
	}
}

func (s *Server) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10000000); err != nil {
		s.InternalError(&w, err)
		return
	}

	if !s.isPostingAllowed(r) {
		http.Error(w, "Not allowed", 403)
		return
	}

	addr, err := s.GetAddress(r)
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	post := &models.Post{}
	post.Title = r.PostFormValue("title")

	file, header, err := r.FormFile("file")
	if err != nil {
		s.InternalError(&w, err)
		return
	}
	defer file.Close()

	if header.Size > s.config.MaxFileSize {
		http.Error(w, fmt.Sprintf("Max file size is %d", s.config.MaxFileSize), 403)
		return
	}
	post.FileSize = header.Size

	for _, ct := range s.config.AllowedContentTypes {
		if ct == header.Header.Get("Content-Type") {
			post.ContentType = ct
		}
	}
	if post.ContentType == "" {
		http.Error(w, fmt.Sprintf("Illegal file type. Allowed types: %s",
			strings.Join(s.config.AllowedContentTypes[:], ", ")), 403)
		return
	}

	if len(post.Title) == 0 && len(header.Filename) < 140 {
		post.Title = header.Filename
	} else if len(post.Title) == 0 {
		post.Title = "."
	}

	var cid string
	for k := range s.ipfs {
		if len(cid) == 0 {
			cid, err = s.ipfs[k].Add(file, ipfs.Pin(s.config.IpfsPin))
			if err != nil {
				s.log.Println("Add error", err)
			}
		} else {
			go s.pinAtShell(k, cid)
		}
	}

	if len(cid) == 0 {
		s.InternalError(&w, errors.New("Failed to add file to IPFS"))
		return
	}

	post.CID = cid
	post.AddressId = addr.Id
	if !s.config.Premoderation {
		post.Published = true
	}

	if err := s.store.Post().Create(post); err != nil {
		s.InternalError(&w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%sres/%d.html", s.config.SitePath, post.Id), 303)
}

func (s *Server) PostViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	post, err := s.store.Post().GetById(vars["id"])
	if err == sql.ErrNoRows {
		http.Error(w, "Not found", 404)
		return
	} else if err != nil {
		s.InternalError(&w, err)
		return
	}
	comments, err := s.store.Comment().GetByPostId(vars["id"])
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	data := struct {
		Post           *models.Post
		Comments       []*models.Comment
		EnableVotes    bool
		EnableComments bool
		Token          string
	}{
		Post:           post,
		Comments:       comments,
		EnableVotes:    s.config.EnableVotes,
		EnableComments: s.config.EnableComments,
		Token:          nosurf.Token(r),
	}

	err = s.template.ExecuteTemplate(w, "post_view.html", data)
	if err != nil {
		s.InternalError(&w, err)
	}
}

func (s *Server) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10000000); err != nil {
		s.InternalError(&w, err)
		return
	}

	addr, err := s.GetAddress(r)
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	comment := &models.Comment{}
	comment.Text = r.PostFormValue("text")

	if comment.PostId, err = strconv.Atoi(r.PostFormValue("post_id")); err != nil {
		s.InternalError(&w, err)
		return
	}

	comment.AddressId = addr.Id

	if err := s.store.Comment().Create(comment); err != nil {
		s.InternalError(&w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("%sres/%d.html#comment%d", s.config.SitePath, comment.PostId, comment.Id), 303)
}

func (s *Server) UpvoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr, err := s.GetAddress(r)
	if err != nil {
		s.InternalError(&w, err)
		return
	}
	vote := &models.Vote{}
	vote.AddressId = addr.Id
	if vote.PostId, err = strconv.Atoi(vars["id"]); err != nil {
		s.InternalError(&w, err)
		return
	}
	if err := s.store.Vote().Upvote(vote); err != nil && err != sql.ErrNoRows {
		s.InternalError(&w, err)
		return
	}
	if vote.Id == 0 {
		w.WriteHeader(403)
	} else {
		w.WriteHeader(201)
	}
}

func (s *Server) DownvoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr, err := s.GetAddress(r)
	if err != nil {
		s.InternalError(&w, err)
		return
	}
	vote := &models.Vote{}
	vote.AddressId = addr.Id
	if vote.PostId, err = strconv.Atoi(vars["id"]); err != nil {
		s.InternalError(&w, err)
		return
	}
	if err := s.store.Vote().Downvote(vote); err != nil && err != sql.ErrNoRows {
		s.InternalError(&w, err)
		return
	}
	if vote.Id == 0 {
		w.WriteHeader(403)
	} else {
		w.WriteHeader(201)
	}
}
