package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Hooks []*Hook

type Hook struct {
	Hook    string `yaml:"hook"`
	Command string `yaml:"command"`
	Auth    *Auth  `yaml:"auth"`
}

func ReadHookFile(path string) (Hooks, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)

	if err != nil {
		return nil, err
	}
	hooks := Hooks{}

	if err := yaml.Unmarshal(body, &hooks); err != nil {
		return nil, err
	}
	return hooks, nil
}
