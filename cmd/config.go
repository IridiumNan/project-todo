/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"path"

	"github.com/IridiumNan/project-todo/internal/config"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Open config file with env editor",
	Long: `This command will open global configuration file.
	It will use env EDITOR first, if not set, it will use vim as default editor`,
	Run: execConfig,
}

func execConfig(cmd *cobra.Command, args []string) {
	configFilePath := path.Join(config.APPConfigDir, config.ConfigFileName)

	err := utils.OpenWithEnvEditor(configFilePath, "vim", utils.ModeEdit)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
