package cmd

import (
	"encoding/json"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/EngineeringKiosk/GermanTechPodcasts/io"
	"github.com/spf13/cobra"
)

// convertJsonToReadmeCmd represents the convertJsonToReadme command
var convertJsonToReadmeCmd = &cobra.Command{
	Use:   "convertJsonToReadme",
	Short: "Converts generated Podcast JSON files into a repository README.md",
	Long: `The source of truth are our YAML files in podcasts/..
Those will be converted and enriched into JSON with the convertYamlToJson command.
To make it human readable, we generate a README.md based on this JSON data.

This command converts the generated JSON information into a human readable README.`,
	RunE: cmdConvertJsonToReadme,
}

func init() {
	rootCmd.AddCommand(convertJsonToReadmeCmd)

	convertJsonToReadmeCmd.Flags().String("json-directory", "", "Directory on where to store the json files")
	convertJsonToReadmeCmd.Flags().String("readme-template", "", "Path to the README template")
	convertJsonToReadmeCmd.Flags().String("readme-output", "", "Path to the README file that will be written")

	convertJsonToReadmeCmd.MarkFlagRequired("json-directory")
	convertJsonToReadmeCmd.MarkFlagRequired("readme-template")
	convertJsonToReadmeCmd.MarkFlagRequired("readme-output")

	convertJsonToReadmeCmd.MarkFlagsRequiredTogether("json-directory", "readme-template", "readme-output")
}

func cmdConvertJsonToReadme(cmd *cobra.Command, args []string) error {
	readmeOutput, err := cmd.Flags().GetString("readme-output")
	if err != nil {
		return err
	}

	readmeTemplate, err := cmd.Flags().GetString("readme-template")
	if err != nil {
		return err
	}

	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	log.Printf("Reading files with extension %s from directory %s", io.JSONExtension, jsonDir)
	jsonFiles, err := io.GetAllFilesFromDirectory(jsonDir, io.JSONExtension)
	if err != nil {
		return err
	}
	log.Printf("%d files found with extension %s in directory %s", len(jsonFiles), io.JSONExtension, jsonDir)

	activePodcasts := make([]*PodcastInformation, 0)
	archivedPodcasts := make([]*PodcastInformation, 0)
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

		if podcastInfo.Archive {
			archivedPodcasts = append(archivedPodcasts, podcastInfo)
		} else {
			activePodcasts = append(activePodcasts, podcastInfo)
		}
	}

	// Sort list by name
	log.Printf("Sorting %d active podcasts by name", len(activePodcasts))
	sort.Slice(activePodcasts, func(i, j int) bool {
		return strings.ToLower(activePodcasts[i].Name) < strings.ToLower(activePodcasts[j].Name)
	})

	log.Printf("Sorting %d archived podcasts by name", len(archivedPodcasts))
	sort.Slice(archivedPodcasts, func(i, j int) bool {
		return strings.ToLower(archivedPodcasts[i].Name) < strings.ToLower(archivedPodcasts[j].Name)
	})

	log.Printf("Read template file %s from disk", readmeTemplate)
	readmeTemplateContent, err := os.ReadFile(readmeTemplate)
	if err != nil {
		return err
	}

	log.Printf("Create target file %s", readmeOutput)
	readmeTarget, err := os.Create(readmeOutput)
	if err != nil {
		return err
	}

	podcasts := struct {
		Active   []*PodcastInformation
		Archived []*PodcastInformation
	}{
		Active:   activePodcasts,
		Archived: archivedPodcasts,
	}

	// Render the template
	log.Printf("Render template and write it into %s ...", readmeOutput)
	t := template.Must(template.New("readme-template").Parse(string(readmeTemplateContent)))
	err = t.Execute(readmeTarget, podcasts)
	if err != nil {
		return err
	}
	log.Printf("Render template and write it into %s ... successful", readmeOutput)

	return nil
}
