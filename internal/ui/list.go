package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"

	"newscl/internal/hn"
)

// identifies which HN feed the list screen is showing
type section int

const (
	sectionPosts section = iota
	sectionJobs
)

func (s section) title() string {
	if s == sectionJobs {
		return "Jobs"
	}
	return "Top Posts"
}

const itemsPerPage = 8

// state for the paginated list screen
type listModel struct {
	section section

	ids     []int
	page    int
	cursor  int
	perPage int

	cache map[int]hn.Item

	loading bool
	spinner spinner.Model
}

func newListModel(sec section) listModel {
	sp := spinner.New(spinner.WithSpinner(spinner.MiniDot))
	return listModel{
		section: sec,
		perPage: itemsPerPage,
		cache:   make(map[int]hn.Item),
		spinner: sp,
	}
}

// switches to the list screen for the given section and kicks off
// the initial id fetch.
func (m Model) enterList(sec section) (Model, tea.Cmd) {
	m.screen = screenList
	m.list = newListModel(sec)
	m.list.loading = true
	return m, tea.Batch(m.list.spinner.Tick, m.fetchIDsCmd(sec))
}

// --- messages ---

type idsLoadedMsg struct {
	section section
	ids     []int
	err     error
}

type itemsLoadedMsg struct {
	section section
	page    int
	items   []hn.Item
	err     error
}

// --- commands ---

const fetchTimeout = 10 * time.Second

func (m Model) fetchIDsCmd(sec section) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()

		var (
			ids []int
			err error
		)
		if sec == sectionJobs {
			ids, err = client.JobStoryIDs(ctx)
		} else {
			ids, err = client.TopStoryIDs(ctx)
		}
		return idsLoadedMsg{section: sec, ids: ids, err: err}
	}
}

// fetches whichever ids on the given page aren't already
// cached. If everything needed is already cached, it returns immediately
// with no items to add.
func (m Model) fetchPageCmd(page int) tea.Cmd {
	client := m.client
	sec := m.list.section
	pageIDs := pageSlice(m.list.ids, page, m.list.perPage)

	var need []int
	for _, id := range pageIDs {
		if _, ok := m.list.cache[id]; !ok {
			need = append(need, id)
		}
	}
	if len(need) == 0 {
		return func() tea.Msg {
			return itemsLoadedMsg{section: sec, page: page}
		}
	}

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()

		items, err := client.Items(ctx, need)
		return itemsLoadedMsg{section: sec, page: page, items: items, err: err}
	}
}

// --- update ---

func (m Model) updateList(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case idsLoadedMsg:
		if msg.section != m.list.section {
			return m, nil // stale response from a screen we've since left
		}
		if msg.err != nil {
			return m.enterError(msg.err, m.list.section), nil
		}
		m.list.ids = msg.ids
		m.list.page = 0
		return m, m.fetchPageCmd(0)

	case itemsLoadedMsg:
		if msg.section != m.list.section || msg.page != m.list.page {
			return m, nil // stale response, e.g. user already changed pages
		}
		if msg.err != nil {
			return m.enterError(msg.err, m.list.section), nil
		}
		for _, item := range msg.items {
			m.list.cache[item.ID] = item
		}
		m.list.loading = false
		m.list.cursor = 0
		return m, nil

	case spinner.TickMsg:
		if !m.list.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.list.spinner, cmd = m.list.spinner.Update(msg)
		return m, cmd

	case tea.KeyPressMsg:
		return m.updateListKey(msg)
	}

	return m, nil
}

func (m Model) updateListKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	if m.list.loading {
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			m.screen = screenHome
		}
		return m, nil
	}

	pageIDs := pageSlice(m.list.ids, m.list.page, m.list.perPage)
	total := totalPages(len(m.list.ids), m.list.perPage)

	switch msg.String() {
	case "esc", "q":
		m.screen = screenHome
	case "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.list.cursor > 0 {
			m.list.cursor--
		}
	case "down", "j":
		if m.list.cursor < len(pageIDs)-1 {
			m.list.cursor++
		}

	case "left", "h":
		if m.list.page > 0 {
			m.list.page--
			m.list.cursor = 0
			m.list.loading = true
			return m, tea.Batch(m.list.spinner.Tick, m.fetchPageCmd(m.list.page))
		}
	case "right", "l":
		if m.list.page < total-1 {
			m.list.page++
			m.list.cursor = 0
			m.list.loading = true
			return m, tea.Batch(m.list.spinner.Tick, m.fetchPageCmd(m.list.page))
		}

	case "enter":
		if m.list.cursor < len(pageIDs) {
			id := pageIDs[m.list.cursor]
			if item, ok := m.list.cache[id]; ok {
				return m.enterModal(item), nil
			}
		}
	}

	return m, nil
}

// --- view ---

func (m Model) viewList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.list.section.title()))
	b.WriteString("\n\n")

	if m.list.loading {
		b.WriteString(m.list.spinner.View() + " " + mutedStyle.Render("Loading…"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("esc back"))
		return boxStyle.Width(60).Render(b.String())
	}

	pageIDs := pageSlice(m.list.ids, m.list.page, m.list.perPage)
	if len(pageIDs) == 0 {
		b.WriteString(mutedStyle.Render("Nothing to show right now."))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("esc back"))
		return boxStyle.Width(60).Render(b.String())
	}

	for i, id := range pageIDs {
		item := m.list.cache[id]
		b.WriteString(renderChoice(itemLine(item), i == m.list.cursor))
		b.WriteString("\n")
	}

	total := totalPages(len(m.list.ids), m.list.perPage)
	b.WriteString("\n")
	b.WriteString(mutedStyle.Render(fmt.Sprintf("Page %d/%d", m.list.page+1, total)))

	var nav []string
	if m.list.page > 0 {
		nav = append(nav, "‹ Previous")
	}
	if m.list.page < total-1 {
		nav = append(nav, "Next ›")
	}
	if len(nav) > 0 {
		b.WriteString("   " + mutedStyle.Render(strings.Join(nav, "    ")))
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("↑/↓ select  •  ←/→ page  •  enter open  •  esc back"))

	return boxStyle.Width(60).Render(b.String())
}

// formats a single row in the list: title, then score/comments (or
// just "[job]" for job postings, which don't have scores or comments).
func itemLine(item hn.Item) string {
	title := truncate(item.Title, 46)
	if item.IsJob() {
		return title + mutedStyle.Render("  [job]")
	}
	meta := fmt.Sprintf("  %d pts, %d comments", item.Score, item.Descendants)
	return title + mutedStyle.Render(meta)
}
