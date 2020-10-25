package build

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jenkins-zh/docker-zh/pkg/common"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

type BuildOptions struct {
	*common.Options

	Username string
	Token    string

	DockerUsername string
	DockerToken    string

	UploadToBintray bool
	UploadToNexus   bool
	NexusUsername   string
	NexusPassword   string
	NexusURL        string
	NexusRepo       string
	NexusFilePath   string

	CleanWAR     bool
	CleanTempDir bool
	CleanImage   bool
	DryRun       bool

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
	cmd.Flags().StringVarP(&buildOptions.Username, "username", "u", "",
		`The username of Bintray API`)
	cmd.Flags().StringVarP(&buildOptions.Token, "token", "t", "",
		`The token of Bintray API`)
	cmd.Flags().StringVarP(&buildOptions.DockerUsername, "docker-username", "", "",
		`The username of docker hub`)
	cmd.Flags().StringVarP(&buildOptions.DockerToken, "docker-token", "", "",
		`The token of docker hub`)
	cmd.Flags().BoolVarP(&buildOptions.CleanWAR, "clean-war", "", true,
		`Clean jenkins.war after uploaded it`)
	cmd.Flags().BoolVarP(&buildOptions.CleanTempDir, "clean-temp", "", true,
		`Clean the temp dir after uploaded it`)
	cmd.Flags().BoolVarP(&buildOptions.CleanImage, "clean-image", "", true,
		`Clean docker image after uploaded it`)
	cmd.Flags().BoolVarP(&buildOptions.UploadToBintray, "upload-to-bintray", "", true,
		`Upload files to Bintray, see also https://bintray.com/. Please attention, the default value might be change in the future`)
	cmd.Flags().BoolVarP(&buildOptions.UploadToNexus, "upload-to-nexus", "", false,
		`Upload files to a Nexus server`)
	cmd.Flags().StringVarP(&buildOptions.NexusURL, "nexus-url", "", "",
		`The URL of a Nexus server`)
	cmd.Flags().StringVarP(&buildOptions.Username, "nexus-username", "", "",
		`The username of a Nexus server`)
	cmd.Flags().StringVarP(&buildOptions.NexusPassword, "nexus-password", "", "",
		`The password of a Nexus server`)
	cmd.Flags().StringVarP(&buildOptions.NexusRepo, "nexus-repo", "", "",
		`The repo of a Nexus server`)
	cmd.Flags().StringVarP(&buildOptions.NexusFilePath, "nexus-filepath", "", "",
		`The filepath of a repo from Nexus`)
	cmd.Flags().BoolVarP(&buildOptions.DryRun, "dry-run", "", false,
		`Do not really do the build action`)
	return
}

func (o *BuildOptions) Run(cmd *cobra.Command, args []string) (err error) {
	// load the config file
	if err = o.ConfigManager.Read(o.ConfigPath); err != nil {
		return
	}

	buildMap := make(map[string]VersionFormula, 0)

	cmd.Println("check the new formulas here")
	formulas := make([]common.CustomFormula, 0)
	for _, formula := range o.ConfigManager.GetFormulas() {
		formulaFile := fmt.Sprintf("formulas/%s.yaml", formula.Name)

		var data []byte
		if data, err = ioutil.ReadFile(formulaFile); err != nil {
			return
		}

		latestMD5 := fmt.Sprintf("%x", md5.Sum(data))
		if formula.MD5 != latestMD5 {
			for _, ver := range o.ConfigManager.GetAllVersions() {
				buildMap[fmt.Sprintf("%s-%s", formula.Name, ver)] = VersionFormula{
					Version: ver,
					Formula: formula,
				}
			}
			// update the md5 of this formula, so we don't need package it again
			formula.MD5 = latestMD5
		}

		formulas = append(formulas, formula)
	}

	cmd.Println("start to build all things")
	cmd.Printf("found new versionFormulas %#v, count %d\n", buildMap, len(buildMap))
	if err = o.dockerLogin(cmd.OutOrStdout()); err != nil {
		return
	}

	for _, versionFormula := range buildMap {
		var path string
		if path, err = o.build(versionFormula, cmd.OutOrStdout()); err != nil {
			cmd.Println("failed in build", versionFormula)
			return
		} else {
			if o.UploadToNexus {
				err = o.uploadToNexus(path, versionFormula.Version, versionFormula.Formula.Name)
				if err != nil {
					fmt.Println("upload error", err)
				}
			}

			if o.UploadToBintray {
				err = o.upload(path, versionFormula.Version, versionFormula.Formula.Name)
				if err != nil {
					fmt.Println("upload error", err)
				}
			}

			if o.CleanWAR {
				if cleanErr := os.RemoveAll(path); cleanErr != nil {
					fmt.Println("clean war file error", err)
				}
			}

			if o.CleanTempDir {
				tempDir := fmt.Sprintf("tmp-%s-%s", versionFormula.Formula.Name, versionFormula.Version)

				if cleanErr := os.RemoveAll(tempDir); cleanErr != nil {
					fmt.Println("clean temp file error", err)
				}
			}

			err = o.dockerPush(versionFormula.Version, versionFormula.Formula.Name, cmd.OutOrStdout())
			if err != nil {
				fmt.Println("docker push error", err)
			}

			if o.CleanImage {
				// clean docker image in order not to take too much of storage
				err = o.dockerRemoveImage(versionFormula.Version, versionFormula.Formula.Name, cmd.OutOrStdout())
				if err != nil {
					fmt.Println("docker remove image error", err)
				}
			}
		}
	}

	// should save the config file for each formula
	// just save it one time because it's easy to achieve
	o.ConfigManager.SetFormulas(formulas)
	if err = o.ConfigManager.Save(); err != nil {
		cmd.PrintErrf("cannot save config file %#v\n", err)
	}
	return
}

func (o *BuildOptions) dockerLogin(writer io.Writer) (err error) {
	if o.DockerUsername == "" || o.DockerToken == "" {
		return
	}

	args := []string{"login", "--username", o.DockerUsername, "--password", o.DockerToken}
	if o.DryRun {
		fmt.Println(args)
	} else {
		cmd := exec.Command("docker", args...)
		cmd.Stderr = writer
		cmd.Stdout = writer
		err = cmd.Run()
	}
	return
}

func (o *BuildOptions) dockerPush(version, formula string, writer io.Writer) (err error) {
	args := []string{"push", fmt.Sprintf("jenkinszh/jenkins-%s:%s", formula, version)}
	if o.DryRun {
		fmt.Println(args)
	} else {
		cmd := exec.Command("docker", args...)
		cmd.Stderr = writer
		cmd.Stdout = writer
		err = cmd.Run()
	}
	return
}

func (o *BuildOptions) dockerRemoveImage(version, formula string, writer io.Writer) (err error) {
	args := []string{"rmi", fmt.Sprintf("jenkinszh/jenkins-%s:%s", formula, version)}
	if o.DryRun {
		fmt.Println(args)
	} else {
		cmd := exec.Command("docker", args...)
		cmd.Stderr = writer
		cmd.Stdout = writer
		err = cmd.Run()
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

	request.SetBasicAuth(o.Username, o.Token)
	if response, err = client.Get(api); err == nil && response.StatusCode == 200 {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			err = json.Unmarshal(data, &files)
		}
	} else if response != nil {
		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			fmt.Println("response", string(data))
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
	request.SetBasicAuth(o.Username, o.Token)
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

func (o *BuildOptions) build(versionFormula VersionFormula, writer io.Writer) (path string, err error) {
	configPathTemplate := fmt.Sprintf("formulas/%s.yaml", versionFormula.Formula.Name)

	var configPath string
	if configPath, err = common.RenderTemplate(configPathTemplate, map[string]string{
		"version": versionFormula.Version,
	}); err != nil {
		return
	}
	defer os.RemoveAll(configPath)

	args := []string{"cwp", "--config-path", configPath,
		"--version", versionFormula.Version, "--tmp-dir",
		fmt.Sprintf("tmp-%s-%s", versionFormula.Formula.Name, versionFormula.Version)}
	if o.DryRun {
		fmt.Println(args)
	} else {
		cmd := exec.Command("jcli", args...)
		cmd.Stderr = writer
		cmd.Stdout = writer
		err = cmd.Run()
		path = fmt.Sprintf("tmp-%s-%s/output/target/jenkins-zh-%s.war", versionFormula.Formula.Name,
			versionFormula.Version, versionFormula.Version)
	}
	return
}

// https://stackoverflow.com/questions/41951455/nexus-3-file-upload-to-hosted-maven-repository
func (o *BuildOptions) uploadToNexus(filepath, version, formula string) (err error) {
	client := &http.Client{}
	api := fmt.Sprintf("%s/repository/%s/%s/%s", o.NexusURL, o.NexusRepo, o.NexusFilePath, filepath)

	fmt.Printf("upload file by %s\n", api)
	if o.DryRun {
		return
	}

	var request *http.Request
	var response *http.Response
	data, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer func() {
		_ = data.Close()
	}()
	if request, err = http.NewRequest("POST", api, data); err != nil {
		return
	}
	request.SetBasicAuth(o.NexusUsername, o.NexusPassword)
	response, err = client.Do(request)

	if response != nil {
		fmt.Println("StatusCode", response.StatusCode, "response", response.Body)

		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			fmt.Println("response", string(data))
		}
	}
	return
}

func (o *BuildOptions) upload(filepath, version, formula string) (err error) {
	client := &http.Client{}
	api := fmt.Sprintf("https://api.bintray.com/content/jenkins-zh/generic/jenkins/%s/jenkins-%s.war", version, formula)

	fmt.Printf("upload file by %s\n", api)
	if o.DryRun {
		return
	}

	var request *http.Request
	var response *http.Response
	data, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer data.Close()
	if request, err = http.NewRequest("PUT", api, data); err != nil {
		return
	}
	request.Header.Add("X-Bintray-Package", "jenkins")
	request.Header.Add("X-Bintray-Version", version)
	request.Header.Add("X-Bintray-Publish", "1")
	request.Header.Add("X-Bintray-Override", "1")
	request.Header.Add("X-Bintray-Explode", "0")

	request.SetBasicAuth(o.Username, o.Token)
	response, err = client.Do(request)

	if response != nil {
		fmt.Println("StatusCode", response.StatusCode, "response", response.Body)

		var data []byte
		if data, err = ioutil.ReadAll(response.Body); err == nil {
			fmt.Println("response", string(data))
		}
	}
	return
}
