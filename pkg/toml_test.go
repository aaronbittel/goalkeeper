package pkg

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestTomlConfig(t *testing.T) {
	name := "my-tasks"
	configPath := "/home/user/goalkeeper"
	csvPath := "/home/user/goalkeeper/csv"
	daily := 120

	input := fmt.Sprintf(`[config]
	name = "%s"
	config_path = "%s"
	csv_path = "%s"

	[goals]
	daily = %d`, name, configPath, csvPath, daily)

	var config TomlDocument

	_, err := toml.Decode(input, &config)
	if err != nil {
		t.Error(err)
	}

	if config.ConfigSection.Filename != name {
		t.Errorf("expected %q, got %q", name, config.ConfigSection.Filename)
	}

	if config.GoalsSection.Daily != daily {
		t.Errorf("expected %d, got %d", daily, config.GoalsSection.Daily)
	}
}
