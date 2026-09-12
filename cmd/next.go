/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/filter"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/runner"
	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: execNext,
}

func execNext(cmd *cobra.Command, args []string) {
	energyStr, err := cmd.Flags().GetString("energy")
	if err != nil {
		slog.Error("while getting energy from flag", "err", err)
		os.Exit(1)
	}

	e := parseEnergyInput(energyStr)
	if e == models.EnergyInvalid {
		slog.Error("while parse energy string", "err", err)
		os.Exit(1)
	}

	f := filter.EnergyFilter(e)

	dataDir, err := utils.SearchBeginCurrentDir()
	if err != nil {
		slog.Error("error when search data dir, please check if you have init your project with todo init command", "err", err)
		os.Exit(1)
	}

	td, err := store.NewTomlDB(dataDir)
	if err != nil {
		slog.Error("error when create a new toml database connection", "err", err)
		os.Exit(1)
	}

	work, err := td.Pop(f)
	if err != nil {
		log.Fatal(err)
	}

	out, file := dataDirWriter(dataDir)
	defer file.Close()

	r := runner.NewWorkTomlRunner(dataDir, work, out)

	err = r.Run(td)
	if err != nil {
		log.Fatal(err)
	}
}

func dataDirWriter(dataDir string) (out io.Writer, file io.Closer) {
	logFilePath := path.Join(dataDir, models.DataDirLogName)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return os.Stdout, nil
	}

	out = io.MultiWriter(logFile, os.Stdout)
	file = logFile

	return
}

func parseEnergyInput(energyStr string) models.Energy {
	parseHint := "parse energy input success"
	switch energyStr {
	case "l", "low", "0":
		slog.Info(parseHint, "energy", "low")
		return models.EnergyLow
	case "m", "medium", "1":
		slog.Info(parseHint, "energy", "medium")
		return models.EnergyMedium
	case "h", "high", "2":
		slog.Info(parseHint, "energy", "high")
		return models.EnergyHigh
	}

	return models.EnergyInvalid
}

func init() {
	rootCmd.AddCommand(nextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nextCmd.PersistentFlags().String("foo", "", "A help for foo")

	energyUsage := `Set the energy status

	l or low or 0		low energy now
	m or medium or 1	medium energy now
	h or high or 2		high energy now
	`
	nextCmd.Flags().StringP("energy", "e", "low", energyUsage)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
