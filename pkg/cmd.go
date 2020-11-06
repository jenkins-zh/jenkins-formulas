package pkg

import (
	"github.com/jenkins-zh/jenkins-formulas/pkg/build"
	"github.com/jenkins-zh/jenkins-formulas/pkg/check"
	"github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cj",
	Short: "Custom Jenkins automatically",
}

// register all commands here
func init() {
	rootCmd.Flags().StringVarP(&commonOptions.ConfigPath, "config-path", "", "config.yaml",
		`the config file path`)

	rootCmd.AddCommand(check.NewCheckCommand(&commonOptions))
	rootCmd.AddCommand(build.NewBuildCommand(&commonOptions))
}

// GetRootCmd only for test purpose
func GetRootCmd() *cobra.Command {
	return rootCmd
}

var commonOptions = common.Options{}

// GetCommonOptions only for test purpose
func GetCommonOptions() *common.Options {
	return &commonOptions
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
