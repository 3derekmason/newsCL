package hn

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient spins up an httptest server that serves canned JSON
// responses for the endpoints this package calls, and returns a Client
// pointed at it.
func newTestClient(t *testing.T, routes map[string]any) *Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := routes[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(body); err != nil {
			t.Fatalf("encoding test response: %v", err)
		}
	}))
	t.Cleanup(srv.Close)

	return &Client{http: srv.Client(), baseURL: srv.URL}
}

func TestTopStoryIDsFiltersOutJobs(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"/topstories.json": []int{1, 2, 3, 4},
		"/jobstories.json": []int{2, 4},
	})

	ids, err := client.TopStoryIDs(context.Background())
	if err != nil {
		t.Fatalf("TopStoryIDs returned error: %v", err)
	}

	want := []int{1, 3}
	if len(ids) != len(want) {
		t.Fatalf("got %v, want %v", ids, want)
	}
	for i, id := range ids {
		if id != want[i] {
			t.Fatalf("got %v, want %v", ids, want)
		}
	}
}

func TestItemDecodesFields(t *testing.T) {
	client := newTestClient(t, map[string]any{
		"/item/42.json": Item{
			ID:          42,
			Type:        "story",
			Title:       "Test Story",
			URL:         "https://example.com",
			Score:       100,
			Descendants: 12,
		},
	})

	item, err := client.Item(context.Background(), 42)
	if err != nil {
		t.Fatalf("Item returned error: %v", err)
	}
	if item.Title != "Test Story" || item.Score != 100 || item.Descendants != 12 {
		t.Fatalf("unexpected item: %+v", item)
	}
	if item.LinkURL() != "https://example.com" {
		t.Fatalf("LinkURL() = %q, want external URL", item.LinkURL())
	}
}

func TestItemLinkURLFallsBackToDiscussion(t *testing.T) {
	item := Item{ID: 7, Type: "story"} // no URL, e.g. an Ask HN post
	want := "https://news.ycombinator.com/item?id=7"
	if got := item.LinkURL(); got != want {
		t.Fatalf("LinkURL() = %q, want %q", got, want)
	}
}
