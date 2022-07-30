package cmd

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/EngineeringKiosk/GermanTechPodcasts/io"
	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// convertYamlToJsonCmd represents the convertYamlToJson command
var convertYamlToJsonCmd = &cobra.Command{
	Use:   "convertYamlToJson",
	Short: "Converts Podcast YAML files into JSON files",
	Long: `The YAML representation of the basic podcast info is more for humans.
For machines we have a JSON format with more information about the podcast available.

This command converts the basic YAML information into JSON format.`,
	RunE: cmdConvertYamlToJson,
}

func init() {
	rootCmd.AddCommand(convertYamlToJsonCmd)

	convertYamlToJsonCmd.Flags().String("yaml-directory", "", "Directory on where to find the yaml files")
	convertYamlToJsonCmd.Flags().String("json-directory", "", "Directory on where to store the json files")

	convertYamlToJsonCmd.MarkFlagRequired("yaml-directory")
	convertYamlToJsonCmd.MarkFlagRequired("json-directory")
	convertYamlToJsonCmd.MarkFlagsRequiredTogether("yaml-directory", "json-directory")
}

func cmdConvertYamlToJson(cmd *cobra.Command, args []string) error {
	yamlDir, err := cmd.Flags().GetString("yaml-directory")
	if err != nil {
		return err
	}

	jsonDir, err := cmd.Flags().GetString("json-directory")
	if err != nil {
		return err
	}

	yamlFiles, err := io.GetAllFilesFromDirectory(yamlDir, io.YAMLExtension)
	if err != nil {
		return err
	}

	jsonFiles, err := io.GetAllFilesFromDirectory(jsonDir, io.JSONExtension)
	if err != nil {
		return err
	}

	// Process every YAML file found and dump it into a JSON
	// file with the same name.
	// If the JSON file already exists, merge it and only update the data
	// that is available in the YAML file.
	for _, f := range yamlFiles {
		absYamlFilePath := filepath.Join(yamlDir, f.Name())
		yamlFileContent, err := ioutil.ReadFile(absYamlFilePath)
		if err != nil {
			return err
		}

		podcastInfo := &PodcastInformation{}
		err = yaml.Unmarshal(yamlFileContent, podcastInfo)
		if err != nil {
			return err
		}

		currentFileExtension := path.Ext(f.Name())
		jsonFileName := f.Name()[0:len(f.Name())-len(currentFileExtension)] + io.JSONExtension
		absJsonFilePath := filepath.Join(jsonDir, jsonFileName)

		// Check if we have a related json file already
		if _, ok := jsonFiles[jsonFileName]; ok {
			// JSON file exists.
			// Read JSON file into Podcast Information structure
			// and overwrite yaml information
			jsonFileContent, err := ioutil.ReadFile(absJsonFilePath)
			if err != nil {
				return err
			}

			podcastJsonInfo := &PodcastInformation{}
			err = json.Unmarshal(jsonFileContent, podcastJsonInfo)
			if err != nil {
				return err
			}

			podcastInfo = mergePodcastInformation(podcastInfo, podcastJsonInfo)
		}

		// Add generated fields
		// TODO Maybe this should be an own command
		podcastInfo.Slug = slug.Make(podcastInfo.Name)

		// Dump data into JSON file
		err = io.WriteJSONFile(absJsonFilePath, podcastInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

// mergePodcastInformation will overwrite a fixed set of
// fields from source into target.
func mergePodcastInformation(source, target *PodcastInformation) *PodcastInformation {
	// Those fields are all fields where
	// the yaml file is the source of truth
	// If the yaml structure will be extended
	//
	// This can be implemented via the reflect package,
	// but for now (and the first prototype) it is good enough.
	target.Name = source.Name
	target.Website = source.Website
	target.PodcastIndexID = source.PodcastIndexID
	target.RSSFeed = source.RSSFeed
	target.Spotify = source.Spotify
	target.Description = source.Description
	target.Tags = source.Tags
	target.WeeklyDownloadsAVG = source.WeeklyDownloadsAVG

	return target
}
