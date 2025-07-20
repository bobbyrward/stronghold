package config

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNotificationsConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `
notifiers:
  - name: "slack"
    type: "webhook"
    url: "https://hooks.slack.com/services/xxx"
  - name: "email"
    type: "smtp"
    url: "smtp://mail.example.com"
`
	var cfg NotificationsConfig
	err := yaml.Unmarshal([]byte(yamlData), &cfg)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	expected := NotificationsConfig{
		Notifiers: []NotificationsConfigNotifier{
			{Name: "slack", Type: "webhook", Url: "https://hooks.slack.com/services/xxx"},
			{Name: "email", Type: "smtp", Url: "smtp://mail.example.com"},
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Unmarshaled config does not match expected.\nGot: %#v\nWant: %#v", cfg, expected)
	}
}
