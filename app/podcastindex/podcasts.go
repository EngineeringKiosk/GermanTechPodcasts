package podcastindex

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// PodcastsService handles communication with the Podcast related
// methods of the PodcastIndex API.
//
// PodcastIndex API docs: https://podcastindex-org.github.io/docs-api/#tag--Podcasts
type PodcastsService service

type Podcast struct {
	Feed PodcastFeed

	// Found reports whether the API returned a feed for the requested id.
	// Once a feed is removed from the PodcastIndex, the API answers with
	// an empty array ("feed": []) instead of a feed object.
	Found bool
}

func (p *Podcast) UnmarshalJSON(data []byte) error {
	var raw struct {
		Feed json.RawMessage `json:"feed,omitempty"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// "feed": [] means: No feeds match this id.
	if len(raw.Feed) == 0 || raw.Feed[0] == '[' {
		return nil
	}

	if err := json.Unmarshal(raw.Feed, &p.Feed); err != nil {
		return err
	}
	p.Found = true

	return nil
}

type PodcastFeed struct {
	ID             int               `json:"id,omitempty"`
	Title          string            `json:"title,omitempty"`
	Image          string            `json:"image,omitempty"`
	Artwork        string            `json:"artwork,omitempty"`
	LastUpdateTime int64             `json:"lastUpdateTime,omitempty"`
	ItunesID       int64             `json:"itunesId,omitempty"`
	EpisodeCount   int               `json:"episodeCount,omitempty"`
	Categories     map[string]string `json:"categories,omitempty"`
}

// Get a single Podcast by Feed ID.
//
// PodcastIndex API docs: https://podcastindex-org.github.io/docs-api/#get-/podcasts/byfeedid
func (s *PodcastsService) GetByFeedID(ctx context.Context, feedID int64) (*Podcast, *http.Response, error) {
	u := fmt.Sprintf("podcasts/byfeedid?id=%d", feedID)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	podcast := new(Podcast)
	resp, err := s.client.Do(ctx, req, podcast)
	if err != nil {
		return nil, resp, err
	}

	return podcast, resp, nil
}
