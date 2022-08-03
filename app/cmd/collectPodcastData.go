package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
		jsonFileContent, err := ioutil.ReadFile(absJsonFilePath)
		if err != nil {
			return err
		}

		podcastInfo := &PodcastInformation{}
		err = json.Unmarshal(jsonFileContent, podcastInfo)
		if err != nil {
			return err
		}

		// Get Podcast info
		log.Printf("Requesting 'Podcasts.GetByFeedID' data from podcast index for feed id %d", int(podcastInfo.PodcastIndexID))
		p, _, err := c.Podcasts.GetByFeedID(context.Background(), int(podcastInfo.PodcastIndexID))
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
		err = downloadFile(p.Feed.Artwork, absImageFilePath)
		if err != nil {
			return err
		}

		podcastInfo.Image = filepath.Join(imageFolder, imageFileName)

		// Get Podcast Episodes info
		log.Printf("Requesting 'Episodes.GetByFeedID' data from podcast index for feed id %d", int(podcastInfo.PodcastIndexID))
		episodes, _, err := c.Episodes.GetByFeedID(context.Background(), int(podcastInfo.PodcastIndexID), 1000)
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

func downloadFile(address, fileName string) error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", defaultUserAgent)
	response, err := client.Do(req)

	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		return errors.New("received non 200 response code")
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
