package ipfd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
)

func (s *Server) PostFeed(w http.ResponseWriter, r *http.Request) {
	posts, err := s.store.Post().GetBoardIndex(0, 50, true)
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	var created time.Time
	if len(posts) > 0 {
		created = posts[0].CreatedAt
	} else {
		created = time.Now()
	}

	feed := &feeds.Feed{
		Title:       s.config.siteName,
		Link:        &feeds.Link{Href: s.config.siteURL},
		Description: "Does it need a description?",
		Created:     created,
	}

	for _, p := range posts {
		link := fmt.Sprintf("%s/res/%d.html", s.config.siteURL, p.Id)
		feedItem := &feeds.Item{
			Title:       p.Title,
			Link:        &feeds.Link{Href: link},
			Id:          link,
			Description: s.getIpfsUrl(p.CID),
			Created:     p.CreatedAt,
		}
		feed.Items = append(feed.Items, feedItem)
	}

	rss, err := feed.ToRss()
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	fmt.Fprint(w, rss)
}

func (s *Server) CommentFeed(w http.ResponseWriter, r *http.Request) {
	comments, err := s.store.Comment().GetAll(0, 50)
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	var created time.Time
	if len(comments) > 0 {
		created = comments[0].CreatedAt
	} else {
		created = time.Now()
	}

	feed := &feeds.Feed{
		Title:       s.config.siteName,
		Link:        &feeds.Link{Href: s.config.siteURL},
		Description: "Comments feed",
		Created:     created,
	}

	for _, p := range comments {
		link := fmt.Sprintf("%s/res/%d.html#comment%d", s.config.siteURL, p.PostId, p.Id)
		feedItem := &feeds.Item{
			Title:       p.Text,
			Link:        &feeds.Link{Href: link},
			Id:          link,
			Description: p.Text,
			Created:     p.CreatedAt,
		}
		feed.Items = append(feed.Items, feedItem)
	}

	rss, err := feed.ToRss()
	if err != nil {
		s.InternalError(&w, err)
		return
	}

	fmt.Fprint(w, rss)
}
