package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

type Entries struct{}
type entry struct {
	Username string `json:"u"`
	Key      string `json:"a"`
}

func getConfigFilename() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, ".retro_config"), nil
}

func (e *Entries) Save(entries []entry) error {
	config, err := getConfigFilename()
	if err != nil {
		return err
	}
	output, err := json.Marshal(entries)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config, output, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (e *Entries) Load() ([]entry, error) {
	config, err := getConfigFilename()
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	entries := []entry{}
	err = json.Unmarshal(contents, &entries)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
