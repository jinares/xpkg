package xtools

import (
	"errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

//BindYaml "gopkg.in/yaml.v2"
func BindYAML(path string, bind interface{}) error {
	if FileExists(path) == false {
		return errors.New("not found")
	}

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, bind)
	return err
}

func ToYAML(data interface{}) (string, error) {
	ds, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(ds), nil
}
func YAMLToJSON(str string, obj interface{}) error {
	return yaml.Unmarshal([]byte(str), obj)
}
