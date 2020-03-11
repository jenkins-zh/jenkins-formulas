package common

import (
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
)

func RenderTemplate(path string, values map[string]string) (result string, err error) {
	var t *template.Template
	if t, err = template.ParseFiles(path); err == nil {
		f,_ := ioutil.TempFile("/tmp", ".yaml")
		err = t.Execute(f, values)

		result = f.Name()
	}
	return
}

type CustomWarPackage struct {
	Bundle Bundle
	BuildSettings BuildSettings
	War CustomWar
	Plugins []Plugin
}

type CASC struct {

}

type Plugin struct {
	GroupId string
	ArtifactId string
	Source Source
}

type Bundle struct {
	GroupId string
	ArtifactId string
	Description string
	Vendor string
}

type BuildSettings struct {
	Docker BuildDockerSetting
}

type BuildDockerSetting struct {
	Base string
	Tag string
	Build bool
}

type CustomWar struct {
	GroupId string
	ArtifactId string
	Source Source
}

type Source struct {
	Version string
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