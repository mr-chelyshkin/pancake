package internal

import (
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
read file.

	parse file and unmarshal it to income struct (as interface).
	return incoming struct with data from file.
*/

func ReadYaml(filename string, object interface{}) (interface{}, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data := object

	if err := yaml.Unmarshal(buf, &data); err != nil {
		return nil, err
	}

	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &object,
		TagName:  "yaml",
	}
	decoder, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return nil, err
	}

	decoder.Decode(data)
	return object, nil
}
