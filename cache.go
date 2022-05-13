package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}

type cached struct {
	data []byte
	etag string
}

type cache struct {
	cache map[string]cached // key: url, value: last known data
	mutex sync.RWMutex
}

func newCache() *cache {
	return &cache{cache: make(map[string]cached), mutex: sync.RWMutex{}}
}

func (c *cache) get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.mutex.RLock()
	last, found := c.cache[url]
	c.mutex.RUnlock()
	if found {
		req.Header.Set("If-None-Match", last.etag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if found && resp.StatusCode == http.StatusNotModified {
		return last.data, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200 OK, got %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.mutex.Lock()
		c.cache[url] = cached{data, etag}
		c.mutex.Unlock()
	}
	return data, nil
}
