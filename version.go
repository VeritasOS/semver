package main

import (
	"fmt"
	"os"

	"github.com/VeritasOS/semver/config"
	"github.com/VeritasOS/tool-upgrade-go"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the version",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// omit early out of date check
		if _, _, err := upgrade.RemoveBackup(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("version      " + config.VERSION)
		fmt.Fprintln(os.Stderr, "")
		upgrade.CheckAndNotifyIfOutOfDate(
			config.ToolName,
			config.VERSION,
			config.RepoBase,
			config.VersionFilePrefix,
			config.VersionStable,
			config.HoursUntilNextUpgradeCheck,
			config.UpgradeCommandName,
		)
		return nil
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
