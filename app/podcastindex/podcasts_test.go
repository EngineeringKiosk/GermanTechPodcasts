package podcastindex

import (
	"encoding/json"
	"testing"
)

// The PodcastIndex API returns an empty array as feed once a feed id is unknown.
// Decoding this must not fail, otherwise a single removed podcast breaks the whole data collection.
func TestPodcastUnmarshalJSON_RemovedFeed(t *testing.T) {
	body := `{"status":"true","query":{"id":"685245"},"feed":[],"description":"No feeds match this id."}`

	p := new(Podcast)
	if err := json.Unmarshal([]byte(body), p); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if p.Found {
		t.Error("Found = true, want false")
	}
	if p.Feed.ID != 0 {
		t.Errorf("Feed.ID = %d, want 0", p.Feed.ID)
	}
}

func TestPodcastUnmarshalJSON_ExistingFeed(t *testing.T) {
	body := `{"status":"true","query":{"id":"685245"},"feed":{"id":685245,"title":"IT@DB","artwork":"https://example.com/cover.png","itunesId":1462447493,"episodeCount":94},"description":"Found matching feed."}`

	p := new(Podcast)
	if err := json.Unmarshal([]byte(body), p); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if !p.Found {
		t.Error("Found = false, want true")
	}
	if p.Feed.ID != 685245 {
		t.Errorf("Feed.ID = %d, want 685245", p.Feed.ID)
	}
	if p.Feed.Title != "IT@DB" {
		t.Errorf("Feed.Title = %q, want %q", p.Feed.Title, "IT@DB")
	}
	if p.Feed.Artwork != "https://example.com/cover.png" {
		t.Errorf("Feed.Artwork = %q, want %q", p.Feed.Artwork, "https://example.com/cover.png")
	}
	if p.Feed.ItunesID != 1462447493 {
		t.Errorf("Feed.ItunesID = %d, want 1462447493", p.Feed.ItunesID)
	}
	if p.Feed.EpisodeCount != 94 {
		t.Errorf("Feed.EpisodeCount = %d, want 94", p.Feed.EpisodeCount)
	}
}
