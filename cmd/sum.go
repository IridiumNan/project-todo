/*
Copyright © 2026 IridiumNan 2930416610@qq.com
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IridiumNan/project-todo/internal/runner"
	"github.com/IridiumNan/project-todo/internal/store"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// sumCmd represents the sum command
var sumCmd = &cobra.Command{
	Use:   "sum",
	Short: "Manage the summary of done works",
	Long: `sum new	 	command select a done work then create new summary file
	sum cd		create a new shell then enter the summary file dir`,
	Run: execSum,
}

func execSum(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("not enough args")
		fmt.Println("todo sum --help for details")
		os.Exit(1)
	}

	switch args[0] {
	case "new":
		execSumNew()
	case "cd":
		execSumCd()
	}
}

func execSumNew() {
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

	td := store.NewTomlDoneDB()

	r := runner.NewDoneWorkRunner(dataDir, os.Stdout, td, nil)

	err = r.Run()
	if err != nil {
		fmt.Printf("error while running sum, err: %s", err)
		os.Exit(1)
	}
}

func execSumCd() {
	fmt.Println("exec sum cd")
}

func init() {
	rootCmd.AddCommand(sumCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sumCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sumCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
