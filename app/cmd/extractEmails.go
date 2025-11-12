package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	libIO "github.com/EngineeringKiosk/GermanTechPodcasts/io"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

type AuthorInfo struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type EmailResult struct {
	PodcastName string       `json:"podcastName"`
	RSSFeed     string       `json:"rssFeed"`
	Emails      []string     `json:"emails"`
	Authors     []AuthorInfo `json:"authors,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// extractEmailsCmd represents the extractEmails command
var extractEmailsCmd = &cobra.Command{
	Use:   "extractEmails",
	Short: "Extract contact email addresses from podcast RSS feeds",
	Long: `This command reads all podcast JSON files and extracts contact email addresses
from their RSS feeds. It looks for email addresses in various RSS fields like 
managingEditor, webMaster, author, and description.`,
	RunE: cmdExtractEmails,
}

func init() {
	rootCmd.AddCommand(extractEmailsCmd)

	extractEmailsCmd.Flags().String("json-directory", "", "Directory containing the podcast JSON files")
	extractEmailsCmd.Flags().String("output", "", "Optional: Output file to write results (if not provided, prints to stdout)")

	extractEmailsCmd.MarkFlagRequired("json-directory")
}

func cmdExtractEmails(cmd *cobra.Command, args []string) error {
	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	outputFile, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}

	log.Printf("Reading files with extension %s from directory %s", libIO.JSONExtension, jsonDir)
	jsonFiles, err := libIO.GetAllFilesFromDirectory(jsonDir, libIO.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), libIO.JSONExtension, jsonDir)

	// Create RSS parser
	client := &http.Client{Timeout: 30 * time.Second}
	fp := gofeed.NewParser()
	fp.Client = client

	var results []EmailResult

	// Process each podcast file
	for _, f := range jsonFiles {
		absJsonFilePath := filepath.Join(jsonDir, f.Name())
		log.Printf("Processing file %s", absJsonFilePath)

		jsonFileContent, err := os.ReadFile(absJsonFilePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", absJsonFilePath, err)
			continue
		}

		podcastInfo := &PodcastInformation{}
		err = json.Unmarshal(jsonFileContent, podcastInfo)
		if err != nil {
			log.Printf("Error unmarshaling JSON from %s: %v", absJsonFilePath, err)
			continue
		}

		if podcastInfo.RSSFeed == "" {
			log.Printf("Skipping %s: no RSS feed URL", podcastInfo.Name)
			continue
		}

		result := processRSSFeed(*podcastInfo, fp)
		results = append(results, result)
	}

	// Output results
	if outputFile != "" {
		return writeResultsToFile(results, outputFile)
	} else {
		return printResults(results)
	}
}

func processRSSFeed(podcast PodcastInformation, fp *gofeed.Parser) EmailResult {
	result := EmailResult{
		PodcastName: podcast.Name,
		RSSFeed:     podcast.RSSFeed,
		Emails:      []string{},
	}

	log.Printf("Parsing RSS feed for %s: %s", podcast.Name, podcast.RSSFeed)

	feed, err := fp.ParseURL(podcast.RSSFeed)
	if err != nil {
		result.Error = err.Error()
		log.Printf("Error parsing RSS feed for %s: %v", podcast.Name, err)
		return result
	}

	// Use maps to avoid duplicates
	emails := make(map[string]bool)
	authors := make(map[string]AuthorInfo) // Use email as key to avoid duplicates

	// Check managingEditor and webMaster fields (available in custom fields for RSS)
	if managingEditor, ok := feed.Custom["managingEditor"]; ok {
		addEmailsFromText(managingEditor, emails)
	}
	if webMaster, ok := feed.Custom["webMaster"]; ok {
		addEmailsFromText(webMaster, emails)
	}

	// Check author fields
	if feed.Author != nil {
		if feed.Author.Email != "" {
			emails[feed.Author.Email] = true
			authors[feed.Author.Email] = AuthorInfo{
				Name:  feed.Author.Name,
				Email: feed.Author.Email,
			}
		}
	}
	for _, author := range feed.Authors {
		if author.Email != "" {
			emails[author.Email] = true
			authors[author.Email] = AuthorInfo{
				Name:  author.Name,
				Email: author.Email,
			}
		}
	}

	// Check description for emails
	if feed.Description != "" {
		addEmailsFromText(feed.Description, emails)
	}

	// Check iTunes fields
	if feed.ITunesExt != nil {
		if feed.ITunesExt.Owner != nil && feed.ITunesExt.Owner.Email != "" {
			emails[feed.ITunesExt.Owner.Email] = true
			authors[feed.ITunesExt.Owner.Email] = AuthorInfo{
				Name:  feed.ITunesExt.Owner.Name,
				Email: feed.ITunesExt.Owner.Email,
			}
		}
		if feed.ITunesExt.Author != "" {
			addEmailsFromText(feed.ITunesExt.Author, emails)
		}
	}

	// Check first 3 episodes for contact info
	episodesToCheck := len(feed.Items)
	if episodesToCheck > 3 {
		episodesToCheck = 3
	}

	for i := 0; i < episodesToCheck; i++ {
		item := feed.Items[i]

		if item.Author != nil && item.Author.Email != "" {
			emails[item.Author.Email] = true
			authors[item.Author.Email] = AuthorInfo{
				Name:  item.Author.Name,
				Email: item.Author.Email,
			}
		}
		for _, author := range item.Authors {
			if author.Email != "" {
				emails[author.Email] = true
				authors[author.Email] = AuthorInfo{
					Name:  author.Name,
					Email: author.Email,
				}
			}
		}
		if item.Description != "" {
			addEmailsFromText(item.Description, emails)
		}
	}

	// Convert maps to slices
	for email := range emails {
		result.Emails = append(result.Emails, email)
	}

	for _, author := range authors {
		result.Authors = append(result.Authors, author)
	}

	if len(result.Emails) > 0 {
		log.Printf("Found %d email(s) and %d author(s) for %s: %v", len(result.Emails), len(result.Authors), podcast.Name, result.Emails)
	} else {
		log.Printf("No emails found for %s", podcast.Name)
	}

	return result
}

func addEmailsFromText(text string, emails map[string]bool) {
	// Regular expression to match email addresses
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := emailRegex.FindAllString(text, -1)

	for _, match := range matches {
		// Clean up the email (remove any trailing punctuation)
		email := strings.Trim(match, ".,;!?")
		emails[email] = true
	}
}

func writeResultsToFile(results []EmailResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		return err
	}

	log.Printf("Results written to %s", filename)
	return nil
}

func printResults(results []EmailResult) error {
	// Print summary to stderr so JSON can be piped cleanly
	log.Printf("Found emails from %d podcasts", len(results))

	foundEmails := 0
	totalEmails := 0
	for _, result := range results {
		if len(result.Emails) > 0 {
			foundEmails++
			totalEmails += len(result.Emails)
		}
	}
	log.Printf("Emails found in %d/%d podcasts (%d total emails)", foundEmails, len(results), totalEmails)

	// Print JSON to stdout
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
