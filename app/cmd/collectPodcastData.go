package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	libIO "github.com/EngineeringKiosk/GermanTechPodcasts/io"
	"github.com/EngineeringKiosk/GermanTechPodcasts/podcastindex"
	"github.com/spf13/cobra"
)

const (
	imageFolder      = "images"
	defaultUserAgent = "EngineeringKiosk-GermanTechPodcasts"
)

// collectPodcastDataCmd represents the collectPodcastData command
var collectPodcastDataCmd = &cobra.Command{
	Use:   "collectPodcastData",
	Short: "Collects additional data per podcast from the PodcastIndex API",
	Long: `We only have basic data about each podcast.
To make the whole project more useful, we aim to collect additional data from the PodcastIndex API.

This command gathers this additional data per Podcast and stores them back into
the genrated JSON files.`,
	RunE: cmdCollectPodcastData,
}

func init() {
	rootCmd.AddCommand(collectPodcastDataCmd)

	collectPodcastDataCmd.Flags().String("json-directory", "", "Directory on where to store the json files")
	collectPodcastDataCmd.Flags().String("api-key", "", "API Key for Podcast Index API")
	collectPodcastDataCmd.Flags().String("api-secret", "", "API Secret for Podcast Index API")

	collectPodcastDataCmd.MarkFlagRequired("json-directory")
	collectPodcastDataCmd.MarkFlagRequired("api-key")
	collectPodcastDataCmd.MarkFlagRequired("api-secret")
}

func cmdCollectPodcastData(cmd *cobra.Command, args []string) error {
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}

	apiSecret, err := cmd.Flags().GetString("api-secret")
	if err != nil {
		return err
	}

	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	log.Printf("Reading files with extension %s from directory %s", libIO.JSONExtension, jsonDir)
	jsonFiles, err := libIO.GetAllFilesFromDirectory(jsonDir, libIO.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), libIO.JSONExtension, jsonDir)

	c := podcastindex.NewClient(nil, apiKey, apiSecret)

	for _, f := range jsonFiles {
		absJsonFilePath := filepath.Join(jsonDir, f.Name())
		log.Printf("Processing file %s", absJsonFilePath)
		jsonFileContent, err := os.ReadFile(absJsonFilePath)
		if err != nil {
			return err
		}

		podcastInfo := &PodcastInformation{}
		err = json.Unmarshal(jsonFileContent, podcastInfo)
		if err != nil {
			return err
		}

		if podcastInfo.PodcastIndexID > 0 {
			// Get Podcast info
			log.Printf("Requesting 'Podcasts.GetByFeedID' data from podcast index for feed id %d", podcastInfo.PodcastIndexID)
			p, _, err := c.Podcasts.GetByFeedID(context.Background(), podcastInfo.PodcastIndexID)
			if err != nil {
				return err
			}

			// Set basic podcast data
			podcastInfo.EpisodeCount = p.Feed.EpisodeCount
			podcastInfo.ItunesID = p.Feed.ItunesID

			// Download cover-image
			imageFileExtension := path.Ext(p.Feed.Artwork)
			// Sometimes we have file extensions like .png?t=1655195362
			// but we only want .png
			if strings.Contains(imageFileExtension, "?") {
				imageFileExtension, _, _ = strings.Cut(imageFileExtension, "?")
			}
			jsonFileExtension := path.Ext(f.Name())
			imageFileName := f.Name()[0:len(f.Name())-len(jsonFileExtension)] + imageFileExtension
			absImageFilePath := filepath.Join(jsonDir, imageFolder, imageFileName)
			log.Printf("Downloading %s into %s", p.Feed.Artwork, absImageFilePath)
			resp, err := downloadFile(p.Feed.Artwork, absImageFilePath)
			if err != nil {
				// This is not very great to check the status code.
				// We do this already in downloadFile.
				// A better solution would be to create a "status code error" type and check
				// this one. But this is now quick and dirty and it works :)
				//
				// If this is any kind of error, exit here.
				if resp == nil || resp.StatusCode == 200 {
					return err
				}

				// If we get a non 200 status code, but we have a target image already
				// (like an old one), it is better to use the old one than failing.
				//
				// Having the latest up to date image is not the highest priority here.
				_, imageExistErr := os.Stat(absImageFilePath)
				if resp.StatusCode != 200 && errors.Is(imageExistErr, os.ErrNotExist) {
					return err
				} else {
					log.Printf("We were not able to download the new image %s", p.Feed.Artwork)
					log.Printf("The pipeline didn't fail, because the previous version %s exists", absImageFilePath)
				}
			}

			podcastInfo.Image = filepath.Join(imageFolder, imageFileName)

			// Get Podcast Episodes info
			log.Printf("Requesting 'Episodes.GetByFeedID' data from podcast index for feed id %d", podcastInfo.PodcastIndexID)
			episodes, _, err := c.Episodes.GetByFeedID(context.Background(), podcastInfo.PodcastIndexID, 1000)
			if err != nil {
				return err
			}

			// Determine time/date of latest episode published
			latestEpisodePublished := int64(0)
			for _, e := range episodes.Items {
				if latestEpisodePublished < e.DatePublished {
					latestEpisodePublished = e.DatePublished
				}
			}
			podcastInfo.LatestEpisodePublished = latestEpisodePublished
		} else {
			log.Printf("Skipping data retrieval from PodcastIndex for %s, because PodcastIndex is %d", absJsonFilePath, podcastInfo.PodcastIndexID)
		}

		// Write the information back to the JSON file
		// Dump data into JSON file
		log.Printf("Write %s to disk ...", absJsonFilePath)
		err = libIO.WriteJSONFile(absJsonFilePath, podcastInfo)
		if err != nil {
			return err
		}
		log.Printf("Write %s to disk ... successful", absJsonFilePath)
	}

	return nil
}

func downloadFile(address, fileName string) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	response, err := client.Do(req)

	if err != nil {
		return response, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return response, fmt.Errorf("received %d as status code, expected 200", response.StatusCode)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return response, err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return response, err
	}

	return response, nil
}
