package cmd

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"

    libIO "github.com/EngineeringKiosk/GermanTechPodcasts/io"
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
)

// tagCount holds tag name and its usage count
type tagCount struct {
    Tag   string
    Count int
}

var tagStatsCmd = &cobra.Command{
    Use:   "tagStats",
    Short: "Show statistics about used tags",
    Long:  `Analyzes all podcast YML files and shows how often each tag is used.`,
    RunE:  cmdTagStats,
}

func init() {
    rootCmd.AddCommand(tagStatsCmd)
    tagStatsCmd.Flags().String("yml-directory", "", "Directory containing the YML files")
    tagStatsCmd.MarkFlagRequired("yml-directory")
}

func cmdTagStats(cmd *cobra.Command, args []string) error {
    ymlDir, err := cmd.Flags().GetString("yml-directory")
    if err != nil {
        return err
    }

    // Read YML files
    ymlFiles, err := libIO.GetAllFilesFromDirectory(ymlDir, ".yml")
    if err != nil {
        return err
    }

    // Map to store tag counts
    tagCounts := make(map[string]int)

    // Process each YML file
    for _, f := range ymlFiles {
        absYmlFilePath := filepath.Join(ymlDir, f.Name())
        ymlFileContent, err := os.ReadFile(absYmlFilePath)
        if err != nil {
            return err
        }

        podcastInfo := &PodcastInformation{}
        err = yaml.Unmarshal(ymlFileContent, podcastInfo)
        if err != nil {
            return err
        }

        // Count tags
        for _, tag := range podcastInfo.Tags {
            tagCounts[tag]++
        }
    }

    // Convert map to slice for sorting
    var tagStats []tagCount
    for tag, count := range tagCounts {
        tagStats = append(tagStats, tagCount{tag, count})
    }

    // Sort by count descending
    sort.Slice(tagStats, func(i, j int) bool {
        return tagStats[i].Count > tagStats[j].Count
    })

    // Output results
    log.Printf("Found %d unique tags\n", len(tagStats))
    fmt.Println("\nTag statistics:")
    fmt.Println("Count\tTag")
    fmt.Println("-----\t---")
    for _, tc := range tagStats {
        fmt.Printf("%5d\t%s\n", tc.Count, tc.Tag)
    }

    return nil
}