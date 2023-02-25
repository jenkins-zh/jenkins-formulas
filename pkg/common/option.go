package common

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Options struct {
	ConfigPath string
}

type CustomConfig struct {
	Formulas []CustomFormula
	LTS      []string
	Weekly   []string
}

type CustomFormula struct {
	Name string
	MD5  string
}

type CustomConfigManager struct {
	CustomConfig *CustomConfig
	ConfigPath   string
}

func (c *CustomConfigManager) Read(path string) (err error) {
	c.ConfigPath = path

	var data []byte
	c.CustomConfig = &CustomConfig{}

	if data, err = ioutil.ReadFile(path); err == nil {
		err = yaml.Unmarshal(data, c.CustomConfig)
	}
	return
}

func (c *CustomConfigManager) GetAllVersions() (allVers []string) {
	allVers = make([]string, 0)
	allVers = append(allVers, c.GetLTSList()...)
	allVers = append(allVers, c.GetWeeklyList()...)
	return
}

func (c *CustomConfigManager) GetLTSList() []string {
	return c.CustomConfig.LTS
}

func (c *CustomConfigManager) GetWeeklyList() []string {
	return c.CustomConfig.Weekly
}

func (c *CustomConfigManager) GetFormulas() []CustomFormula {
	return c.CustomConfig.Formulas
}

func (c *CustomConfigManager) SetFormulas(formulas []*CustomFormula) {
	newFormulas := make([]CustomFormula, len(formulas))
	for i := range formulas {
		newFormulas[i] = *formulas[i]
	}
	c.CustomConfig.Formulas = newFormulas
}

func (c *CustomConfigManager) Save() (err error) {
	if c.ConfigPath == "" {
		err = fmt.Errorf("no config file path provide")
		return
	}

	var data []byte
	if data, err = yaml.Marshal(c.CustomConfig); err == nil {
		err = ioutil.WriteFile(c.ConfigPath, data, 0664)
	}
	return
}

func (c *CustomConfigManager) HasTLS(lts string) (exists bool, err error) {
	for _, item := range c.CustomConfig.LTS {
		if item == lts {
			exists = true
			break
		}
	}
	return
}

func (c *CustomConfigManager) AddTLS(lts string) (err error) {
	c.CustomConfig.LTS = append(c.CustomConfig.LTS, lts)

	if cut := len(c.CustomConfig.LTS) - 10; cut > 0 {
		c.CustomConfig.LTS = c.CustomConfig.LTS[cut:]
	}
	return
}

func (c *CustomConfigManager) HasWeekly(weekly string) (exists bool, err error) {
	for _, item := range c.CustomConfig.Weekly {
		if item == weekly {
			exists = true
			break
		}
	}
	return
}

func (c *CustomConfigManager) AddWeekly(weekly string) (err error) {
	c.CustomConfig.Weekly = append(c.CustomConfig.Weekly, weekly)

	if cut := len(c.CustomConfig.Weekly) - 10; cut > 0 {
		c.CustomConfig.Weekly = c.CustomConfig.Weekly[cut:]
	}
	return
}
