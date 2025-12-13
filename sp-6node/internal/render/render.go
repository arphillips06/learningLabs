package render

import (
	"bytes"
	"learninglabs/sp-6node/internal/model"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

func LoadYAML(path string) (model.DeviceConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return model.DeviceConfig{}, err
	}
	var cfg model.DeviceConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return model.DeviceConfig{}, err
	}
	return cfg, nil
}

func RenderTemplate(templatePath string, cfg model.DeviceConfig) (string, error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, cfg); err != nil {
		return "", err
	}
	return buf.String(), nil
}
