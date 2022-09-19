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

	podcasts := make([]*PodcastInformation, 0)
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

		podcasts = append(podcasts, podcastInfo)
	}

	log.Printf("Sorting %d podcasts by name", len(podcasts))
	// Sort list by name
	sort.Slice(podcasts, func(i, j int) bool {
		return strings.ToLower(podcasts[i].Name) < strings.ToLower(podcasts[j].Name)
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
