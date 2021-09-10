package utils

import "gopkg.in/yaml.v2"

func YamlLoad(model interface{}, data []byte) error  {
	if err := yaml.Unmarshal(data, model); err != nil {
		return err
	}
	return nil
}