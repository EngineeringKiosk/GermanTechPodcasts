package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "GermanTechPodcasts",
	Short: "Helper tooling for a curated list of Tech Podcasts in German",
	Long: `A curated list of Tech Podcasts in German.
To help with this, we automate as much as possible.
This is the helper tooling around the podcast collection.

More information at https://github.com/EngineeringKiosk/GermanTechPodcasts`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
