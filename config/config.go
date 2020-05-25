package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/wailsapp/wails"
)

var settings Settings

func GetSettings() Settings {
	return settings
}

type Config struct {
	Logger *wails.CustomLogger
}

type Settings struct {
	Profile  string              `json:"profile"`
	Profiles map[string]Profile  `json:"profiles"`
	Layout   string              `json:"layout"`
	Layouts  map[string]string   `json:"layouts"`
	Orders   map[string][]string `json:"orders"`
}

func NewConfig() *Config {
	return &Config{}
}

func getConfigFilename() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, ".retro_config"), nil
}

func (e *Config) WailsInit(runtime *wails.Runtime) error {
	e.Logger = runtime.Log.New("config")
	e.Logger.Info("Config Initialization")

	name, err := getConfigFilename()
	if err != nil {
		return err
	}
	// check for existance of config file
	_, err = os.Stat(name)
	if os.IsNotExist(err) {
		e.Logger.Warn("No config file found, using defaults")
		settings = Settings{
			Profile:  "",
			Profiles: map[string]Profile{},
			Layout:   "",
			Layouts:  map[string]string{},
		}
		return nil
	}
	// if we got past the existance, assume we can read it
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		e.Logger.Errorf("Could not open config file: %w", err)
		return err
	}
	err = json.Unmarshal(contents, &settings)
	if err != nil {
		return err
	}

	e.Logger.Infof("Successfully loaded config")
	return nil
}

func (e *Config) save() error {
	config, err := getConfigFilename()
	if err != nil {
		return err
	}
	output, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(config, output, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
