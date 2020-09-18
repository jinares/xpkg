package xtools

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

//BindYaml "gopkg.in/yaml.v2"
func BindYaml(path string, bind interface{}) error {
	if FileExists(path) == false {
		return errors.New("not found")
	}

	yamlFile, err := ioutil.ReadFile(path)

	err = yaml.Unmarshal(yamlFile, bind)
	return err
}
