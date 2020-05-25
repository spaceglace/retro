package config

import (
	"fmt"
)

type Profile struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (e *Config) GetProfiles() map[string]Profile {
	return settings.Profiles
}

func (e *Config) AddProfile(name, key string) error {
	settings.Profiles[name] = Profile{
		Name: name,
		Key:  key,
	}
	return e.SetActiveProfile(name)
}

func (e *Config) RemoveProfile(name string) error {
	delete(settings.Profiles, name)
	return e.save()
}

func (e *Config) GetActiveProfile() string {
	return settings.Profile
}

func (e *Config) SetActiveProfile(name string) error {
	// ensure profile exists
	if _, ok := settings.Profiles[name]; !ok {
		return fmt.Errorf("Profile '%s' does not exist", name)
	}
	settings.Profile = name
	return e.save()
}
