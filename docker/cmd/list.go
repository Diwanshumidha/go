package cmd

import (
	"dock/service"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List containers",
	Long: `List containers.

It can be used to list all containers or filter them by a specific name.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.ListContainers()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
