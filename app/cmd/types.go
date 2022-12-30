package cmd

import (
	"strings"
	"time"
)

const (
	// Number of months until a podcast is marked as dead (aka archived)
	TimeToArchive = -6
)

type PodcastInformation struct {
	Name                   string    `yaml:"name" json:"name"`
	Slug                   string    `json:"slug"`
	Website                string    `yaml:"website" json:"website"`
	PodcastIndexID         int64     `yaml:"podcastIndexID" json:"podcastIndexID"`
	RSSFeed                string    `yaml:"rssFeed" json:"rssFeed"` // TODO Should be better a url.URL
	Spotify                string    `yaml:"spotify" json:"spotify"` // TODO Should be better a url.URL
	Description            string    `yaml:"description" json:"description"`
	Tags                   []string  `yaml:"tags" json:"tags"`
	WeeklyDownloadsAVG     Statistic `yaml:"weekly_downloads_avg" json:"weekly_downloads_avg"`
	EpisodeCount           int       `json:"episodeCount"`
	LatestEpisodePublished int64     `json:"latestEpisodePublished"`
	Archive                bool      `json:"archive"`
	ItunesID               int64     `json:"itunesID"`
	Image                  string    `json:"image"`
}

type Statistic struct {
	Value   int64  `yaml:"value" json:"value"`
	Updated string `yaml:"updated" json:"updated"`
}

func (p PodcastInformation) GetHumanReadableDate() string {
	t := time.Unix(p.LatestEpisodePublished, 0)
	s := t.Format("Monday, 02 January 2006")

	return s
}

func (p PodcastInformation) TagsAsList() string {
	s := ""
	if len(p.Tags) > 0 {
		s = strings.Join(p.Tags, ", ")
	}

	return s
}

// GetLastEpisodeStatus calculates a traffic light status
// on when the last episode was published.
//
// Legend:
//
//	🔴 Last Episode published > 6 months ago
//	🟡 Last Episode published something between 2 months and 6 months ago
//	🟢 Last Episode published within today and last 2 month
func (p PodcastInformation) GetLastEpisodeStatus() string {
	t := time.Unix(p.LatestEpisodePublished, 0)

	sixMonth := time.Now().AddDate(0, -6, 0)
	twoMonth := time.Now().AddDate(0, -2, 0)

	if t.Before(sixMonth) {
		return "🔴"
	}

	if t.After(sixMonth) && t.Before(twoMonth) {
		return "🟡"
	}

	return "🟢"
}
