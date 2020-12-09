package core

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//ReadYaml reads a YAML file
func ReadYaml(path string, out interface{}) (err error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("ReadYaml - Invalid file %s: %v", path, err)
		return
	}

	err = yaml.Unmarshal(d, out)
	if err != nil {
		logrus.Errorf("ReadYaml - Invalid file %s: %v", path, err)
		return
	}
	return
}

//WriteYaml writes a yaml file
func WriteYaml(path string, in interface{}) (err error) {
	d, err := yaml.Marshal(in)
	if err != nil {
		logrus.Errorf("WriteYaml - Cannot marshal to file %s: %v", path, err)
		return
	}
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		logrus.Errorf("WriteYaml - Cannot save file %s: %v", path, err)
		return
	}
	return
}
