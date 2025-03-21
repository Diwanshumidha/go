package cmd

import (
	"dock/internal/constant"
	"dock/internal/dockerClient"
	"dock/service"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   constant.New.Use,
	Short: constant.New.Short,
	Long:  constant.New.Long,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("New command can be used to create a new container with existing templates.")
		fmt.Println("")
		fmt.Println("To use this command, run dock new <template-name>")
		fmt.Println("")

		// TODO: Implement select options
		fmt.Printf("Templates:\n")
		fmt.Printf("- %s\n", "postgres")
		fmt.Printf("- %s\n", "redis")

	},
}

var newPostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Create a new Postgres container",
	Long: `Create a new Postgres container.

It can be used to create a new Postgres container.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Postgres command can be used to create a new Postgres container.")
		if err := dockerClient.MakeSureDockerIsRunning(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		service.NewPostgresCmd()
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.AddCommand(newPostgresCmd)
}
