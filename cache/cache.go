// Package cache implements a cache for HTTP responses, using the server's ETag to determine freshness.

package cache

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// client is an HTTP client with a timeout.
var client = &http.Client{Timeout: 10 * time.Second}

// cached represents previously-fetched data and its ETag.
type cached struct {
	data []byte
	etag string
}

// Cache stores the data for a requested URL.
type Cache struct {
	cache map[string]cached
	mutex sync.RWMutex
}

// New builds a new cache and initializes the map.
func New() *Cache {
	return &Cache{cache: make(map[string]cached), mutex: sync.RWMutex{}}
}

// Get fetches a URL, using the previously-cached data if the server returns 304 Not Modified.
func (c *Cache) Get(url string) (data []byte, found bool, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}

	c.mutex.RLock()
	last, found := c.cache[url]
	c.mutex.RUnlock()
	if found {
		req.Header.Set("If-None-Match", last.etag)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}
	if found && resp.StatusCode == http.StatusNotModified {
		return last.data, true, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("expected status code 200 OK, got %d", resp.StatusCode)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.mutex.Lock()
		c.cache[url] = cached{data, etag}
		c.mutex.Unlock()
	}
	return data, true, nil
}
