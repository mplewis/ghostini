package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/justincampbell/timeago"
)

type host struct {
	apiUrl     string
	contentKey string
}

type Post struct {
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	PublishedAt  time.Time `json:"published_at"`
	PublishedAgo string
	Excerpt      string `json:"excerpt"`
	ReadingTime  int    `json:"reading_time"`
}

type postsResp struct {
	Posts []Post `json:"posts"`
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

func getPosts(c *cache, h host, page int) (postsResp, error) {
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
	data, err := c.get(url)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(data, &resp)
	for i := range resp.Posts {
		resp.Posts[i].PublishedAgo = timeago.FromTime(resp.Posts[i].PublishedAt)
	}
	return resp, err
}
