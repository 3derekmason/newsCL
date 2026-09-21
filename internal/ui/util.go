package ui

import (
	"net/url"
	"strings"
)

// shortens s to at most max runes, appending an ellipsis if it had
// to cut anything off. Used to keep long HN titles from wrapping awkwardly.
func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return string(r[:max])
	}
	return string(r[:max-1]) + "…"
}

// extracts just the host from a URL, e.g. "https://foo.com/bar" ->
// "foo.com". Used as a small visual hint of where a link goes, the way a
// browser's address bar or a news aggregator would show it. Returns "" if
// the item has no external URL (e.g. an Ask HN post).
func domain(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Hostname(), "www.")
}

// returns the slice of ids belonging to the given zero-indexed
// page. Returns nil if page is past the end of ids.
func pageSlice(ids []int, page, perPage int) []int {
	start := page * perPage
	if start >= len(ids) {
		return nil
	}
	end := start + perPage
	if end > len(ids) {
		end = len(ids)
	}
	return ids[start:end]
}

// returns how many pages of perPage items it takes to cover n
// items. Always at least 1
func totalPages(n, perPage int) int {
	if n == 0 {
		return 1
	}
	pages := (n + perPage - 1) / perPage
	if pages == 0 {
		return 1
	}
	return pages
}
