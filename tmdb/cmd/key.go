/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"tmdb/internal/key"

	"github.com/spf13/cobra"
)

// keyCmd represents the key command
var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage your TMDB API key",
	Long: `Manage your TMDB API key stored in the system keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := key.GetAPIKey()
		if err != nil {
			fmt.Println("❗You don't have an API key set. Use `tmdb key set <key>` to set it.")
			return
		}

		fmt.Println("✅ Your API key is set.")
		fmt.Println("Use `tmdb` to get the latest movies.")
	},
}


var keySetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set your API key",
	Long: `Set your API key. This will be stored in your system's keyring.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		argsKey := args[0]
		if err := key.ValidateAPIKey(argsKey); err != nil {
			fmt.Printf("❗ Invalid API key: %s\n", argsKey)
			return
		}
		if err := key.SaveAPIKey(argsKey); err != nil {
			fmt.Printf("❗ Error saving API key: %s\n", argsKey)
			return
		}

		fmt.Printf("✅ API key saved: %s\n", argsKey)
	},
}

var keyDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete your API key",
	Long: `Delete your API key. This will remove it from your system's keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := key.DeleteAPIKey(); err != nil {
			fmt.Printf("❗ Error deleting API key\n")
			return
		}

		fmt.Printf("✅ API key deleted\n")
	},
}

var keyListCmd = &cobra.Command{
	Use:   "get",
	Short: "List your API key",
	Long: `List your API key. This will list all API key stored in your system's keyring.`,
	Run: func(cmd *cobra.Command, args []string) {
		key, err := key.GetAPIKey()
		if err != nil {
			fmt.Printf("❗ Error listing API keys: %s\n", err)
			return
		}
		fmt.Printf("✅ API key: %s\n", key)
	},
}

func init() {
	keyCmd.AddCommand(keySetCmd,keyListCmd, keyDeleteCmd, keySetCmd)
	rootCmd.AddCommand(keyCmd)
}
