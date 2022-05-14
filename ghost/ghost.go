// Package ghost implements a client for the Ghost Content API: https://ghost.org/docs/content-api/

package ghost

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/justincampbell/timeago"
	"github.com/mplewis/ghostini/cache"
	"github.com/mplewis/ghostini/parse"
)

var perPage = fmt.Sprintf("%d", parse.Int(os.Getenv("PER_PAGE"), 10))

// Host represents the target Ghost instance to fetch data from.
type Host struct {
	APIURL     string
	ContentKey string
}

// PostMeta is the metadata for a post, returned when fetching several posts.
type PostMeta struct {
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Excerpt      string    `json:"excerpt"`
	ReadingTime  int       `json:"reading_time"`
	PublishedAt  time.Time `json:"published_at"`
	PublishedAgo string
}

// Post is the full data for a post when fetched individually.
type Post struct {
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	HTML         string    `json:"html"`
	ReadingTime  int       `json:"reading_time"`
	CreatedAt    time.Time `json:"created_at"`
	PublishedAt  time.Time `json:"published_at"`
	PublishedAgo string
	UpdatedAt    time.Time `json:"updated_at"`
	Updated      bool
	UpdatedAgo   string
}

// PostsResp is the response from the Ghost API for several posts.
type PostsResp struct {
	Posts []PostMeta `json:"posts"`
	Meta  struct {
		Pagination struct {
			Page  int `json:"page"`
			Limit int `json:"limit"`
			Pages int `json:"pages"`
			Total int `json:"total"`
			Next  int `json:"next"`
			Prev  int `json:"prev"`
		} `json:"pagination"`
	} `json:"meta"`
}

// PostResp is the response from the Ghost API for a single post.
type PostResp struct {
	Posts []Post `json:"posts"`
}

// GetPosts fetches a page of posts from Ghost.
func GetPosts(c *cache.Cache, h Host, page int) (PostsResp, error) {
	v := url.Values{}
	v.Set("page", fmt.Sprintf("%d", page))
	v.Set("limit", perPage)
	v.Set("key", h.ContentKey)
	v.Set("fields", "title,slug,published_at,excerpt,reading_time")
	v.Set("filter", "visibility:public")
	// HACK: Ghost won't return reading time without plaintext, or excerpt without HTML
	// https://github.com/TryGhost/Ghost/issues/10396#issuecomment-918849637
	v.Set("formats", "plaintext,html")
	url := fmt.Sprintf("%s/ghost/api/v4/content/posts/?%s", h.APIURL, v.Encode())

	var resp PostsResp
	data, _, err := c.Get(url)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(data, &resp)
	for i := range resp.Posts {
		resp.Posts[i].PublishedAgo = timeago.FromTime(resp.Posts[i].PublishedAt)
	}
	return resp, err
}

// GetPost fetches a single post from Ghost.
func GetPost(c *cache.Cache, h Host, slug string) (r PostResp, found bool, err error) {
	v := url.Values{}
	v.Set("key", h.ContentKey)
	url := fmt.Sprintf("%s/ghost/api/v4/content/posts/slug/%s?%s", h.APIURL, slug, v.Encode())

	var resp PostResp
	data, found, err := c.Get(url)
	if err != nil || !found {
		return resp, found, err
	}
	err = json.Unmarshal(data, &resp)
	for i := range resp.Posts {
		resp.Posts[i].PublishedAgo = timeago.FromTime(resp.Posts[i].PublishedAt)
		resp.Posts[i].UpdatedAgo = timeago.FromTime(resp.Posts[i].UpdatedAt)
		resp.Posts[i].Updated = resp.Posts[i].UpdatedAt.After(resp.Posts[i].CreatedAt)
	}
	return resp, found, err
}
