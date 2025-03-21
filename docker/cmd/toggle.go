package cmd

import (
	"dock/internal/constant"
	"dock/service"

	"github.com/spf13/cobra"
)

var toggleCmd = &cobra.Command{
	Use:   constant.Toggle.Use,
	Short: constant.Toggle.Short,
	Long:  constant.Toggle.Long,
	Run: func(cmd *cobra.Command, args []string) {
		service.ToggleContainerUsingPrompt()
	},
}

func init() {
	rootCmd.AddCommand(toggleCmd)
}
