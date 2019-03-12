package main

import (
	"errors"

	"github.com/spf13/cobra"
)

var upgradeForce *bool

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades semver",
	Long:  "Upgrades semver",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Error: Upgrade functionality is unavailable as we transition to publishing this as open source")
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeForce = upgradeCmd.Flags().BoolP("force", "f", false, "force upgrade to latest")
}
