package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"
)

// Page is what build_complete is told about each page of a site.
type Page struct {
	URL         string              `json:"url"`
	Kind        string              `json:"kind"`
	Title       string              `json:"title"`
	Type        string              `json:"type"`
	Excerpt     string              `json:"excerpt"`
	Text        string              `json:"text"`
	Params      map[string]any      `json:"params"`
	Taxonomies  map[string][]string `json:"taxonomies"`
	PublishedAt *time.Time          `json:"published_at"`
}

// Input is build_complete's input.
type Input struct {
	Settings struct {
		FullText *bool `json:"full_text"`
	} `json:"settings"`
	Pages []Page `json:"pages"`
}

// Entry is one page in the index the browser searches.
type Entry struct {
	URL   string     `json:"url"`
	Title string     `json:"title"`
	Date  *time.Time `json:"date,omitempty"`
	Tags  []string   `json:"tags,omitempty"`
	Text  string     `json:"text"`
}

// Index lists the pages a reader can search for: every post and page, newest
// first, but the ones whose front matter says search: false.
func Index(in Input) ([]byte, error) {
	full := in.Settings.FullText == nil || *in.Settings.FullText
	entries := []Entry{}
	for _, p := range in.Pages {
		if p.Kind != "single" || off(p.Params["search"]) {
			continue
		}
		e := Entry{URL: p.URL, Title: p.Title, Date: p.PublishedAt, Text: p.Excerpt}
		if full && p.Text != "" {
			e.Text = p.Text
		}
		for _, name := range slices.Sorted(maps.Keys(p.Taxonomies)) {
			for _, term := range p.Taxonomies[name] {
				if !slices.Contains(e.Tags, term) {
					e.Tags = append(e.Tags, term)
				}
			}
		}
		entries = append(entries, e)
	}
	slices.SortStableFunc(entries, func(a, b Entry) int {
		switch {
		case a.Date == nil && b.Date != nil:
			return 1
		case a.Date != nil && b.Date == nil:
			return -1
		case a.Date != nil && !a.Date.Equal(*b.Date):
			return b.Date.Compare(*a.Date)
		}
		return cmp.Compare(a.URL, b.URL)
	})
	return json.Marshal(entries)
}

// off reads a front matter switch the way an author writes it.
func off(v any) bool {
	if v == nil {
		return false
	}
	switch strings.ToLower(fmt.Sprint(v)) {
	case "false", "no", "off":
		return true
	}
	return false
}
