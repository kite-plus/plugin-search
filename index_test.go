package main

import (
	"encoding/json"
	"slices"
	"testing"
	"time"
)

func TestIndexHoldsPostsAndPagesNewestFirst(t *testing.T) {
	day := func(d int) *time.Time {
		at := time.Date(2026, 9, d, 8, 0, 0, 0, time.UTC)
		return &at
	}
	full := false
	in := Input{Pages: []Page{
		{URL: "/", Kind: "home", Title: "Notes"},
		{URL: "/posts/old/", Kind: "single", Title: "Old", Text: "Old words", PublishedAt: day(1),
			Taxonomies: map[string][]string{"tags": {"go", "wasm"}, "categories": {"tech", "go"}}},
		{URL: "/posts/new/", Kind: "single", Title: "New", Excerpt: "New", Text: "New words", PublishedAt: day(20)},
		{URL: "/about/", Kind: "single", Title: "About", Text: "About me"},
		{URL: "/posts/hidden/", Kind: "single", Title: "Hidden", Params: map[string]any{"search": "no"}},
		{URL: "/tags/go/", Kind: "term", Title: "go"},
	}}

	data, err := Index(in)
	if err != nil {
		t.Fatal(err)
	}
	var got []Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	var urls []string
	for _, e := range got {
		urls = append(urls, e.URL)
	}
	if want := []string{"/posts/new/", "/posts/old/", "/about/"}; !slices.Equal(urls, want) {
		t.Fatalf("urls = %q, want %q", urls, want)
	}
	if !slices.Equal(got[1].Tags, []string{"tech", "go", "wasm"}) || got[1].Text != "Old words" {
		t.Errorf("old = %+v", got[1])
	}

	in.Settings.FullText = &full
	data, err = Index(in)
	if err != nil {
		t.Fatal(err)
	}
	var short []Entry
	if err := json.Unmarshal(data, &short); err != nil {
		t.Fatal(err)
	}
	if short[0].Text != "New" || short[1].Text != "" {
		t.Errorf("without full text, got %q and %q, want the excerpts", short[0].Text, short[1].Text)
	}
}
