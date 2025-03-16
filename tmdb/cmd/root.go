package cmd

import (
	"fmt"
	"os"

	Config "tmdb/internal/config"
	"tmdb/internal/key"
	"tmdb/internal/tmdb"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "tmdb",
	Short: "A utility to get the latest movies from TMDB",
	Long: `Tmdb is a CLI tool that allows you to get the latest movies from TMDB.

It can be used to get the latest movies from TMDB and save them to a file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := Config.Config{}
		fType, _ := cmd.Flags().GetString("type")

		if err := config.SetFolderType(fType); err != nil {
			return fmt.Errorf("error: %w", err)
		}

		fmt.Printf("Fetching '%s' movies...\n", config.Type)

		apiKey, err := key.GetAPIKey()
		if err != nil {
			return fmt.Errorf("API key not set. Use `tmdb key set <key>` to set it.")
		}

		movies, err := tmdb.GetMovies(config.Type, apiKey)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}

		tmdb.DisplayMovies(movies)

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("type", "popular", "Type of movies to get (options: popular, top, upcoming, playing)")
}
