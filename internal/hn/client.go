// client for the public Hacker News Firebase API
// (https://github.com/HackerNews/API)
package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

const defaultBaseURL = "https://hacker-news.firebaseio.com/v0"

// fetches data from the Hacker News API.
type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient() *Client {
	return &Client{
		http:    &http.Client{},
		baseURL: defaultBaseURL,
	}
}

// performs a GET request against url and decodes the JSON response
// body into out.
func (c *Client) getJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("hn: building request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("hn: request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("hn: unexpected status %d from %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("hn: decoding response from %s: %w", url, err)
	}

	return nil
}

// returns up to 200 ids from /v0/jobstories.json
func (c *Client) JobStoryIDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := c.getJSON(ctx, c.baseURL+"/jobstories.json", &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// returns up to 500 ids from /v0/topstories.json, with any job
// postings filtered out.
//
// The topstories feed is a mixed bag of stories, Ask HN, Show HN, and job
// posts all together. Rather than fetching every single item's full body
// just to check its "type" field (500 extra HTTP requests!), we fetch the
// jobstories id list in parallel and drop any ids that shows up there
// too. That keeps "View Posts" and "View Jobs" from ever overlapping.
func (c *Client) TopStoryIDs(ctx context.Context) ([]int, error) {
	var (
		topIDs, jobIDs []int
		topErr, jobErr error
		wg             sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		topErr = c.getJSON(ctx, c.baseURL+"/topstories.json", &topIDs)
	}()
	go func() {
		defer wg.Done()
		jobErr = c.getJSON(ctx, c.baseURL+"/jobstories.json", &jobIDs)
	}()
	wg.Wait()

	if topErr != nil {
		return nil, topErr
	}
	if jobErr != nil {
		return topIDs, nil
	}

	jobSet := make(map[int]struct{}, len(jobIDs))
	for _, id := range jobIDs {
		jobSet[id] = struct{}{}
	}

	filtered := make([]int, 0, len(topIDs))
	for _, id := range topIDs {
		if _, isJob := jobSet[id]; !isJob {
			filtered = append(filtered, id)
		}
	}

	return filtered, nil
}

// fetches a single item by id from /v0/item/{id}.json
func (c *Client) Item(ctx context.Context, id int) (Item, error) {
	var item Item
	url := fmt.Sprintf("%s/item/%d.json", c.baseURL, id)
	if err := c.getJSON(ctx, url, &item); err != nil {
		return Item{}, err
	}
	return item, nil
}

// fetches multiple items concurrently
func (c *Client) Items(ctx context.Context, ids []int) ([]Item, error) {
	items := make([]Item, len(ids))
	errs := make([]error, len(ids))

	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		go func(i, id int) {
			defer wg.Done()
			item, err := c.Item(ctx, id)
			items[i] = item
			errs[i] = err
		}(i, id)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}
