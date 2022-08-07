package podcastindex

import (
	"context"
	"fmt"
	"net/http"
)

// EpisodesService handles communication with the episodes related
// methods of the PodcastIndex API.
//
// PodcastIndex API docs: https://podcastindex-org.github.io/docs-api/#tag--Episodes
type EpisodesService service

type Episodes struct {
	Items []Episode `json:"items,omitempty"`
}

type Episode struct {
	ID            int    `json:"id,omitempty"`
	Title         string `json:"title,omitempty"`
	Link          string `json:"link,omitempty"`
	Description   string `json:"description,omitempty"`
	DatePublished int64  `json:"datePublished,omitempty"`
	Duration      int    `json:"duration,omitempty"`
	Episode       int    `json:"episode,omitempty"`
	Season        int    `json:"season,omitempty"`
	Image         string `json:"image,omitempty"`
}

// Get all episodes by Feed ID.
//
// PodcastIndex API docs: https://podcastindex-org.github.io/docs-api/#get-/episodes/byfeedid
func (s *EpisodesService) GetByFeedID(ctx context.Context, feedID int64, max int) (*Episodes, *http.Response, error) {
	u := fmt.Sprintf("episodes/byfeedid?id=%d&max=%d", feedID, max)
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	episodes := new(Episodes)
	resp, err := s.client.Do(ctx, req, episodes)
	if err != nil {
		return nil, resp, err
	}

	return episodes, resp, nil
}
