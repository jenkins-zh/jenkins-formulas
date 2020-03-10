package build

import (
	"fmt"
	"github.com/jenkins-zh/docker-zh/pkg/common"
	"github.com/spf13/cobra"
)

// NewBuildCommand build the custom Jenkins
func NewBuildCommand(commonOpts *common.Options) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "build",
		Short: "build the custom Jenkins",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not support")
		},
	}
	return
}