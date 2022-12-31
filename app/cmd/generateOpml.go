package cmd

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/EngineeringKiosk/GermanTechPodcasts/io"
	"github.com/spf13/cobra"
)

// generateOpmlCmd represents the generateOpml command
var generateOpmlCmd = &cobra.Command{
	Use:   "generateOpml",
	Short: "Generates an OPML file based on the Podcast JSON files",
	Long: `OPML is an established standard for interop between outliners and RSS readers.
This format is also used by many Podcatchers to enable an import of Podcasts.

This command generates an OPML file based on the JSON information.
See http://opml.org/ for more.`,
	RunE: generateOpml,
}

func init() {
	rootCmd.AddCommand(generateOpmlCmd)

	generateOpmlCmd.Flags().String("json-directory", "", "Directory on where to store the json files")
	generateOpmlCmd.Flags().String("opml-output", "", "Path to the README file that will be written")

	generateOpmlCmd.MarkFlagRequired("json-directory")
	generateOpmlCmd.MarkFlagRequired("opml-output")

	generateOpmlCmd.MarkFlagsRequiredTogether("json-directory", "opml-output")
}

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    OPMLHead `xml:"head"`
	Body    OPMLBody `xml:"body"`
}

type OPMLHead struct {
	Title       string `xml:"title"`
	DateCreated string `xml:"dateCreated"`
	OwnerName   string `xml:"ownerName"`
	OwnerEmail  string `xml:"ownerEmail"`
}

type OPMLBody struct {
	Outlines []OPMLOutline `xml:"outline"`
}

type OPMLOutline struct {
	Title   string `xml:"title,attr"`
	Text    string `xml:"text,attr"`
	Type    string `xml:"type,attr"`
	XMLUrl  string `xml:"xmlUrl,attr"`
	HTMLUrl string `xml:"htmlUrl,attr"`
}

func generateOpml(cmd *cobra.Command, args []string) error {
	opmlOutput, err := cmd.Flags().GetString("opml-output")
	if err != nil {
		return err
	}

	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	t := time.Now()
	o := OPML{
		Version: "2.0",
		Head: OPMLHead{
			Title:       "Deutschsprachige Tech Podcasts",
			DateCreated: t.UTC().Format(time.RFC1123),
			OwnerName:   "Engineering Kiosk",
			OwnerEmail:  "stehtisch@engineeringkiosk.dev",
		},
		Body: OPMLBody{
			Outlines: make([]OPMLOutline, 0),
		},
	}

	log.Printf("Reading files with extension %s from directory %s", io.JSONExtension, jsonDir)
	jsonFiles, err := io.GetAllFilesFromDirectory(jsonDir, io.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), io.JSONExtension, jsonDir)

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

		podcastOutline := OPMLOutline{
			Title:   podcastInfo.Name,
			Text:    podcastInfo.Name,
			Type:    "rss",
			XMLUrl:  podcastInfo.RSSFeed,
			HTMLUrl: podcastInfo.Website,
		}
		o.Body.Outlines = append(o.Body.Outlines, podcastOutline)
	}

	// Sort list by name
	log.Printf("Sorting %d active podcasts by name", len(o.Body.Outlines))
	sort.Slice(o.Body.Outlines, func(i, j int) bool {
		return strings.ToLower(o.Body.Outlines[i].Title) < strings.ToLower(o.Body.Outlines[j].Title)
	})

	log.Printf("Generating OPML file and write it into %s ... ", opmlOutput)
	out, err := xml.MarshalIndent(o, " ", "  ")
	if err != nil {
		return err
	}

	opmlContent := xml.Header + string(out)
	err = os.WriteFile(opmlOutput, []byte(opmlContent), 0644)
	if err != nil {
		return err
	}
	log.Printf("Generating OPML file and write it into %s ... successful", opmlOutput)

	return nil
}
