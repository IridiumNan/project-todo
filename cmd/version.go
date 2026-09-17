/*
Copyright © 2026 IridiumNan 2930416610@qq.com
*/
package cmd

import (
	"fmt"

	"github.com/IridiumNan/project-todo/internal/models"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Show version of this cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("project-todo version %s\n", models.VERSION)
		fmt.Println("For more information, see https://github.com/IridiumNan/project-todo")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
