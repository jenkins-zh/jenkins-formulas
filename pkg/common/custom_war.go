package common

import (
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
)

func RenderTemplate(path string, values map[string]string) (result string, err error) {
	var t *template.Template
	if t, err = template.ParseFiles(path); err == nil {
		f, _ := ioutil.TempFile("/tmp", ".yaml")
		err = t.Execute(f, values)

		result = f.Name()
	}
	return
}

type CustomWarPackage struct {
	Bundle        Bundle        `yaml:"bundle"`
	BuildSettings BuildSettings `yaml:"buildSettings"`
	War           CustomWar     `yaml:"war"`
	Plugins       []Plugin      `yaml:"plugins"`
}

type CASC struct {
}

type Plugin struct {
	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
	Source     Source `yaml:"source"`
}

type Bundle struct {
	GroupId     string `yaml:"groupId"`
	ArtifactId  string `yaml:"artifactId"`
	Description string `yaml:"description"`
	Vendor      string `yaml:"vendor"`
}

type BuildSettings struct {
	Docker BuildDockerSetting `yaml:"docker"`
}

type BuildDockerSetting struct {
	Base  string `yaml:"base"`
	Tag   string `yaml:"tag"`
	Build bool   `yaml:"build"`
}

type CustomWar struct {
	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
	Source     Source `yaml:"source"`
}

type Source struct {
	Version string `yaml:"version"`
}

func ReadCustomWarConfig(path string) (cwp *CustomWarPackage, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err == nil {
		err = yaml.Unmarshal(data, cwp)
	}
	return
}

func SetCustomWarConfigVersion(version string, cwp *CustomWarPackage) {
	cwp.War.Source.Version = version
}

func SaveCustomWarConfig(cwp *CustomWarPackage, path string) (err error) {
	return
}
