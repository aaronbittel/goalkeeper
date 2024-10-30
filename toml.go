package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	TOML_CONFIG_NAME = "config.toml"
)

type ConfigSection struct {
	Filename    string `toml:"name"`
	GoalMinutes int    `toml:"goal"`
}

type TomlDocument struct {
	ConfigSection ConfigSection `toml:"config"`
}

func DefaultTomlConfig() TomlDocument {
	return TomlDocument{
		ConfigSection: ConfigSection{
			Filename: "my-tasks.csv",
		},
	}
}

func loadTomlConfig() TomlDocument {
	var tomlDoc TomlDocument
	_, err := toml.DecodeFile(TOML_CONFIG_NAME, &tomlDoc)

	if err != nil {
		if os.IsExist(err) {
			log.Fatalf("could not parse config.toml: %v", err)
		}
		tomlDoc := DefaultTomlConfig()
		err := createTomlFile(tomlDoc)
		if err != nil {
			log.Fatal(err)
		}
	}

	return tomlDoc
}

func createTomlFile(config TomlDocument) error {
	f, err := os.Create(TOML_CONFIG_NAME)
	if err != nil {
		return fmt.Errorf("error creating config.toml: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding toml config (%v): %v", config, err)
	}

	log.Printf("created default toml config (%s)\n", TOML_CONFIG_NAME)
	return nil
}
