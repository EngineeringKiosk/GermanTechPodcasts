package podcastindex

import (
	"context"
	"fmt"
	"net/http"
)

// PodcastsService handles communication with the Podcast related
// methods of the PodcastIndex API.
//
// PodcastIndex API docs: https://podcastindex-org.github.io/docs-api/#tag--Podcasts
type PodcastsService service

type Podcast struct {
	Feed PodcastFeed `json:"feed,omitempty"`
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
