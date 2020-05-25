package config

import "fmt"

func (e *Config) GetLayouts() map[string]string {
	if len(settings.Layouts) == 0 {
		e.AddLayout(
			"default",
			`{"name":"default","width":300,"height":600,"auto":false,"interval":30,"widgets":[]}`,
		)
	}
	return settings.Layouts
}

func (e *Config) AddLayout(name string, serialized string) error {
	// check if the name is already in use
	if _, ok := settings.Layouts[name]; ok {
		return fmt.Errorf("Name '%s' is already in use", name)
	}

	settings.Layouts[name] = serialized
	return e.save()
}

func (e *Config) RemoveLayout(name string) error {
	delete(settings.Layouts, name)
	return e.save()
}

func (e *Config) UpdateLayout(name string, serialized string) error {
	settings.Layouts[name] = serialized
	return e.save()
}

func (e *Config) GetActiveLayout() string {
	if settings.Layout == "" {
		return "default"
	}
	return settings.Layout
}

func (e *Config) SetActiveLayout(name string) error {
	// ensure layout exists
	if _, ok := settings.Layouts[name]; !ok {
		return fmt.Errorf("Layout '%s' does not exist", name)
	}
	settings.Layout = name
	return e.save()
}
