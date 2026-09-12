/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "A command opening the log file with specific program, default use less",
	Long: `log <command>     will open log file with this command.
	You can run 

	todo log less
	todo log tail
	todo log vim
	etc...`,
	Run: execLog,
}

func execLog(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		slog.Warn("the command not found, use editor to open log file")
		err := openLogWithEditor()
		if err != nil {
			slog.Error("Fail to exec log command", "err", err.Error())
			os.Exit(1)
		}

		return
	}

	program := strings.Trim(args[0], "\n \t")
	err := openLogWithProgram(program)
	if err != nil {
		slog.Error("Fail to exec log command", "err", err)
		os.Exit(1)
	}
}

func openLogWithEditor() error {
	err := utils.OpenWithEnvEditor(utils.GetDefaultLogPath(), "vim")
	if err != nil {
		return fmt.Errorf("error when open log file with env editor, err: %s", err)
	}

	return nil
}

func openLogWithProgram(program string) error {
	err := utils.OpenWithProgram(program, utils.GetDefaultLogPath())
	if err != nil {
		return fmt.Errorf("error when open log file with program, program: %s, err: %s", program, err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
