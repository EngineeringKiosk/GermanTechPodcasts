package cmd

type PodcastInformation struct {
	Name           string   `yaml:"name" json:"name"`
	Slug           string   `json:"slug"`
	Website        string   `yaml:"website" json:"website"`
	PodcastIndexID int64    `yaml:"podcastIndexID" json:"podcastIndexID"`
	RSSFeed        string   `yaml:"rssFeed" json:"rssFeed"` // TODO Should be better a url.URL
	Spotify        string   `yaml:"spotify" json:"spotify"` // TODO Should be better a url.URL
	Description    string   `yaml:"description" json:"description"`
	EpisodeCount   int      `json:"episodeCount"`
	LastUpdateTime int64    `json:"lastUpdateTime"`
	ItunesID       int64    `json:"itunesID"`
	Categories     []string `json:"categories"`
	Image          string   `json:"image"`
}
