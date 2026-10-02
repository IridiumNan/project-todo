/*
Copyright © 2026 IridiumNan 2930416610@qq.com
*/
package cmd

import (
	"log"
	"os"
	"path"

	"github.com/IridiumNan/project-todo/internal/config"
	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/IridiumNan/project-todo/internal/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A terminal tool for store and manage project todo list by Energy",
	Long: `The project-todo will store your todo works then support dependencies tree build.
	You provide your energy status then it will pop an avialable work suitable for now.
	And resolve the dependencies.
	
	If run without any arguments, it will open the global todo doc`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: execRoot,
}

// execRoot open the global todo list recored file
// Just open it wile the configed editor
func execRoot(cmd *cobra.Command, args []string) {
	editor := config.GlobalConf.GlobalTodoEditor

	dataDir := models.GetGlobalDataHome()

	filePath := path.Join(dataDir, models.DataGlobalTodoDoc)

	err := utils.OpenWithProgram(editor, filePath, utils.NoFlag)
	if err != nil {
		log.Fatal(err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.project-todo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
