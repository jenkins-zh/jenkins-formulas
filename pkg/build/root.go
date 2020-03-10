package build

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/docker-zh/pkg/common"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

type BuildOptions struct {
	*common.Options

	Username string
	Password string

	ConfigManager *common.CustomConfigManager
}

// NewBuildCommand build the custom Jenkins
func NewBuildCommand(commonOpts *common.Options) (cmd *cobra.Command) {
	buildOptions := BuildOptions{
		Options: commonOpts,
		ConfigManager: &common.CustomConfigManager{},
	}
	cmd = &cobra.Command{
		Use:   "build",
		Short: "build the custom Jenkins",
		RunE: buildOptions.Run,
	}
	return
}

func (o *BuildOptions) Run(cmd *cobra.Command, args []string) (err error) {
	// load the config file
	if err = o.ConfigManager.Read(o.ConfigPath); err != nil {
		return
	}
	//curl -ulinuxsuren:$TOKEN  https://api.bintray.com/packages/jenkins-zh/generic/jenkins/versions/2.204.5
	//{"name":"2.204.5","desc":null,"package":"jenkins","repo":"generic","owner":"jenkins-zh",
	// "labels":["chinese","jenkins"],"published":true,"attribute_names":[],"created":"2020-03-08T14:37:13.554Z",
	// "updated":"2020-03-08T14:37:47.753Z","released":"","ordinal":3.0,"github_release_notes_file":null,
	// "github_use_tag_release_notes":false,"vcs_tag":null,"rating_count":0}
	//{"message":"Version ':2.204.5' was not found"}


	buildMap := make(map[string]VersionFormula, 0)
	for _, ver := range o.ConfigManager.GetAllVersions() {
		var files []BintrayFile
		if files, err = o.getVersionFiles(ver); err != nil {
			cmd.Println("get version files error", ver, err)
			files = make([]BintrayFile, 0)
		}

		for _, formula := range o.ConfigManager.GetFormulas() {
			found := false
			for _, file := range files {
				if file.Name == fmt.Sprintf("jenkins-%s.war", formula) {
					found = true
					break
				}
			}

			if !found {
				buildMap[fmt.Sprintf("%s-%s", formula.Name, ver)] = VersionFormula{
					Version: ver,
					Formula: formula,
				}
			}
		}
	}

	cmd.Println("check the new formulas here")
	for _, formula := range o.ConfigManager.GetFormulas() {
		formulaFile := fmt.Sprintf("formulas/%s.yaml", formula.Name)

		var data []byte
		if data, err = ioutil.ReadFile(formulaFile); err != nil {
			return
		}

		if formula.MD5 != fmt.Sprintf("%x", md5.Sum(data)) {
			for _, ver := range o.ConfigManager.GetAllVersions() {
				buildMap[fmt.Sprintf("%s-%s", formula.Name, ver)] = VersionFormula{
					Version: ver,
					Formula: formula,
				}
			}
		}
	}

	cmd.Println("start to build all things")
	cmd.Println("found new versionFormulas", len(buildMap))
	for _, versionFormula := range buildMap {
		if err = o.build(versionFormula); err != nil {
			cmd.Println("failed in build", versionFormula)
			return
		}
	}
	return
}

type VersionFormula struct {
	Version string
	Formula common.CustomFormula
}

type BintrayVersion struct {
	Name string
	Desc string
	Package string
	Repo string
	Owner string
	Labels []string
	Published bool
	AttributeNames []string
	Created string
	Updated string
	Released string
	Message string
}
//[{"name":"jenkins-zh.war","path":"jenkins/2.204.5/jenkins-zh.war","repo":"generic","package":"jenkins",
// "version":"2.204.5","owner":"jenkins-zh","created":"2020-03-08T14:37:18.206Z","size":65771513,
// "sha1":"6f03423fdee9fadd332736cfc4e037bed8853c90","sha256":"526ec96de6b32cda2ef16a63ef9c92a32485254737f76606c2105546059204b5"},
// {"name":"jenkins-pipeline-zh.war","path":"jenkins/2.204.5/jenkins-pipeline-zh.war","repo":"generic",
// "package":"jenkins","version":"2.204.5","owner":"jenkins-zh","created":"2020-03-08T14:37:36.865Z",
// "size":93238314,"sha1":"71283205f0d18cb3824288e0960a7f0c5ade8a8b","sha256":"7e4bb1d54db077fb73f917cf134d01fb70f4c2e2d7c5b199c8244804fe2c047f"}]

type BintrayFile struct {
	Name string
	Path string
	Repo string
	Package string
	Version string
	Owner string
	Created string
	Size int64
	Sha1 string
	Sha256 string
}

func (o *BuildOptions) getVersionFiles(version string) (files []BintrayFile, err error){
	client := &http.Client{}
	api := fmt.Sprintf("https://api.bintray.com/packages/jenkins-zh/generic/jenkins/versions/%s/files", version)

	var request *http.Request
	var response *http.Response

	if request, err = http.NewRequest("GET", api, nil); err != nil {
		return
	}

	files = make([]BintrayFile, 0)
	request.SetBasicAuth(o.Username, o.Password)
	if response, err = client.Get(api); err == nil {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			err = json.Unmarshal(data, files)
		}
	}
	return
}

func (o *BuildOptions) checkVersion(version string) (exists bool, err error) {
	client := &http.Client{}
	api := fmt.Sprintf("https://api.bintray.com/packages/jenkins-zh/generic/jenkins/versions/%s", version)

	var request *http.Request
	var response *http.Response

	if request, err = http.NewRequest("GET", api, nil); err != nil {
		return
	}

	bintrayVersion := &BintrayVersion{}
	request.SetBasicAuth(o.Username, o.Password)
	if response, err = client.Get(api); err == nil {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			err = json.Unmarshal(data, bintrayVersion)
		}
	}

	if err == nil {
		exists = bintrayVersion.Published
	}
	return
}

func (o *BuildOptions) build(versionFormula VersionFormula) (err error) {
	var jcli string
	if jcli, err = exec.LookPath("jcli"); err != nil {
		return
	}

	err = syscall.Exec(jcli, []string{"jcli", " cwp", "--config-path",
		fmt.Sprintf("formulas/%s.yaml", versionFormula.Formula),
		"--version", versionFormula.Version, "--tmp-dir", " tmp"}, os.Environ())
	return
}
