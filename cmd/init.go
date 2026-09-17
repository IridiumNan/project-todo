/*
Copyright © 2026 IridiumNan 2930416610@qq.com
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/IridiumNan/project-todo/internal/cache"
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
	You should use the command on the root dir of your project
	USAGE: 

	todo init
	This method init the data dir on current dir

	todo init /path/to/dir
	This method init with current data dir`,
	Run: execInit,
}

func execInit(cmd *cobra.Command, args []string) {
	// NOTE: The cobra args not contains init itself
	// If run `program init`, len(args) will be 0
	logFile := utils.SetGlobalLogger()
	defer logFile.Close()

	var dataDir string
	var err error

	if len(args) == 0 {
		dataDir, err = store.InitDataDirOnCurrentDir()
	} else {
		err = store.InitDataDir(args[0])
		dataDir = args[0]
	}
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("\nData dir has existed on current dir, ignore init command\n\n")
		} else {
			slog.Error("error when init data dir", "err", err)
		}
	}

	dirCache, err := cache.DefaultCache()
	if err != nil {
		slog.Error("creating a new dir cache", "err", err)
	}

	err = dirCache.PushNewDir(dataDir)
	if err != nil {
		slog.Error("pushing current data dir into cache file", "err", err)
		return
	}

	slog.Info("caching current data dir to cache file", "data_dir", dataDir, "cache_path", cache.DefaultDirCacheFilePath)
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
