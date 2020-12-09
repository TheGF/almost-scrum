package core

import (
	"encoding/json"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

//ReadJSON reads a JSON file
func ReadJSON(path string, out interface{}) (err error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("ReadJSON - Invalid file %s: %v", path, err)
		return
	}

	err = json.Unmarshal(d, out)
	if err != nil {
		logrus.Errorf("ReadJSON - Invalid file %s: %v", path, err)
		return
	}
	return
}

//WriteJSON writes a JSON file
func WriteJSON(path string, in interface{}) (err error) {
	d, err := json.Marshal(in)
	if err != nil {
		logrus.Errorf("WriteJSON - Cannot marshal to file %s: %v", path, err)
		return
	}
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		logrus.Errorf("WriteJSON - Cannot save file %s: %v", path, err)
		return
	}
	return
}
