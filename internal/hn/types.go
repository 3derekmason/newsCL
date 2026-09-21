package hn

import "fmt"

// Item represents a single Hacker News item
// https://github.com/HackerNews/API
type Item struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Text        string `json:"text"`
	Score       int    `json:"score"`
	Descendants int    `json:"descendants"`
	Kids        []int  `json:"kids"`
	Dead        bool   `json:"dead"`
	Deleted     bool   `json:"deleted"`
}

// reports whether this item is a job posting
func (i Item) IsJob() bool {
	return i.Type == "job"
}

// returns the item's own page on news.ycombinator.com
func (i Item) DiscussionURL() string {
	return fmt.Sprintf("https://news.ycombinator.com/item?id=%d", i.ID)
}

// returns the URL "View Post"/"View Job" should open: the item's
// external link if it has one, or its own HN page otherwise
func (i Item) LinkURL() string {
	if i.URL != "" {
		return i.URL
	}
	return i.DiscussionURL()
}
