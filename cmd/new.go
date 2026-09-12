/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/IridiumNan/project-todo/internal/builder"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new work then push it to todo list with statusTODO",
	Long: `todo new

	This command will create a new work then push it into dataDir.
	You will edit it's configuration then it will parse markdown text automatically.
	It will search for dataDir which named as .project-todo first, if not found, it will throw error`,
	Run: execNew,
}

func execNew(cmd *cobra.Command, args []string) {
	logFile := utils.SetGlobalLogger()
	defer logFile.Close()

	pwd, err := os.Getwd()
	if err != nil {
		slog.Error("New: error when get current dir, exiting", "err", err)
		os.Exit(1)
	}

	dataDir, err := utils.SearchDataDir(pwd)

	if os.IsNotExist(err) {
		slog.Error("Data Dir not found, please run init command on your project root dir first")
		os.Exit(1)
	}
	slog.Info("found data dir, use it", "path", dataDir)

	wb := builder.NewWorkTomlBuilder(dataDir)

	err = wb.Build()
	if err != nil {
		slog.Error("error when build a new work", "err", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
