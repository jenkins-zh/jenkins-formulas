package check

import (
	"github.com/jenkins-zh/jenkins-formulas/pkg/common"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
)

type CheckOptions struct {
	*common.Options

	ConfigManager *common.CustomConfigManager
}

// NewCheckCommand check the version command
func NewCheckCommand(commonOpts *common.Options) (cmd *cobra.Command) {
	checkOptions := CheckOptions{
		Options:       commonOpts,
		ConfigManager: &common.CustomConfigManager{},
	}
	cmd = &cobra.Command{
		Use:   "check",
		Short: "check update of the Jenkins version",
		RunE:  checkOptions.Run,
	}
	return
}

func (o *CheckOptions) Run(cmd *cobra.Command, args []string) (err error) {
	client := &http.Client{}

	// load the config file
	if err = o.ConfigManager.Read(o.ConfigPath); err != nil {
		return
	}

	var response *http.Response
	if response, err = client.Get("https://updates.jenkins.io/stable/latestCore.txt"); err == nil && response.StatusCode == 200 {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			ver := string(data)

			cmd.Printf("get the latest lts version is %s\n", ver)
			err = o.UpdateLTSVersion(ver)
		}
	}

	if err != nil {
		return
	}

	if response, err = client.Get("https://updates.jenkins.io/latestCore.txt"); err == nil && response.StatusCode == 200 {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			ver := string(data)

			cmd.Printf("get the latest weekly version is %s\n", ver)
			err = o.UpdateWeeklyVersion(ver)
		}
	}

	if err == nil {
		err = o.ConfigManager.Save()
	}

	return
}

func (o *CheckOptions) UpdateLTSVersion(version string) (err error) {
	var ok bool
	if ok, err = o.ConfigManager.HasTLS(version); err == nil && !ok {
		err = o.ConfigManager.AddTLS(version)
	}
	return
}

func (o *CheckOptions) UpdateWeeklyVersion(version string) (err error) {
	var ok bool
	if ok, err = o.ConfigManager.HasWeekly(version); err == nil && !ok {
		err = o.ConfigManager.AddWeekly(version)
	}
	return
}
