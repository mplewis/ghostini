package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

type host struct {
	apiUrl     string
	contentKey string
}

type postsResp struct {
	Posts []struct {
		Title     string    `json:"title"`
		Slug      string    `json:"slug"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"posts"`
	Meta struct {
		Pagination struct {
			Page  int         `json:"page"`
			Limit int         `json:"limit"`
			Pages int         `json:"pages"`
			Total int         `json:"total"`
			Next  int         `json:"next"`
			Prev  interface{} `json:"prev"`
		} `json:"pagination"`
	} `json:"meta"`
}

func getPosts(c *cache, h host, page int) (postsResp, error) {
	v := url.Values{}
	v.Set("page", fmt.Sprintf("%d", page))
	v.Set("key", h.contentKey)
	v.Set("fields", "title,slug,created_at")
	v.Set("filter", "visibility:public")
	url := fmt.Sprintf("%s/ghost/api/v4/content/posts/?%s", h.apiUrl, v.Encode())

	var resp postsResp
	data, err := c.get(url)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(data, &resp)
	return resp, err
}
