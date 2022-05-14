package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/justincampbell/timeago"
	"github.com/mplewis/ghostini/cache"
)

type host struct {
	apiUrl     string
	contentKey string
}

type PostMeta struct {
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	Excerpt      string    `json:"excerpt"`
	ReadingTime  int       `json:"reading_time"`
	PublishedAt  time.Time `json:"published_at"`
	PublishedAgo string
}

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

type postsResp struct {
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

type postResp struct {
	Posts []Post `json:"posts"`
}

func getPosts(c *cache.Cache, h host, page int) (postsResp, error) {
	v := url.Values{}
	v.Set("page", fmt.Sprintf("%d", page))
	v.Set("limit", "10")
	v.Set("key", h.contentKey)
	v.Set("fields", "title,slug,published_at,excerpt,reading_time")
	v.Set("filter", "visibility:public")
	// HACK: Ghost won't return reading time without plaintext, or excerpt without HTML
	// https://github.com/TryGhost/Ghost/issues/10396#issuecomment-918849637
	v.Set("formats", "plaintext,html")
	url := fmt.Sprintf("%s/ghost/api/v4/content/posts/?%s", h.apiUrl, v.Encode())

	var resp postsResp
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

func getPost(c *cache.Cache, h host, slug string) (r postResp, found bool, err error) {
	v := url.Values{}
	v.Set("key", h.contentKey)
	url := fmt.Sprintf("%s/ghost/api/v4/content/posts/slug/%s?%s", h.apiUrl, slug, v.Encode())

	var resp postResp
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
