/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"

	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// , initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a project-todo on current dir",
	Long: `This command will create a new dir for todo data storage
	It will generate a new dir .project-todo.
	You should use the command on the root dir of your project`,
	Run: execInit,
}

func execInit(cmd *cobra.Command, args []string) {
	logFile := utils.SetGlobalLogger()
	defer logFile.Close()

	dataDir, err := store.InitDataDir()
	if err != nil {
		slog.Error("exec init failed", "err", err, "data_dir_path", dataDir)
		return
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
