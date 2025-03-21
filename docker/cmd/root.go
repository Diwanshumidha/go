package cmd

import (
	"dock/internal/constant"
	"dock/service"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   constant.Root.Use,
	Short: constant.Root.Short,
	Long:  constant.Root.Long,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		id, _ := cmd.Flags().GetString("id")

		if id != "" {
			service.ToggleContainerByID(id)
			return
		}

		if name != "" {
			service.ToggleContainerByName(name)
			return
		}

		service.ToggleContainerUsingPrompt()
	},
}

func init() {
	rootCmd.Flags().StringP("name", "n", "", "Name of the container")
	rootCmd.Flags().StringP("id", "i", "", "ID of the container")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
