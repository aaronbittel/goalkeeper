package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	DEFAULT_CONFIG_NAME = "config.toml"
	DEFAULT_CSV_NAME    = "my-tasks.csv"
)

type ConfigSection struct {
	Filename string `toml:"name"`
}

type GoalsSection struct {
	Daily int `toml:"daily"`
}

type TomlDocument struct {
	ConfigSection ConfigSection `toml:"config"`
	GoalsSection  GoalsSection  `toml:"goals"`
}

func DefaultTomlConfig() TomlDocument {
	return TomlDocument{
		ConfigSection: ConfigSection{
			Filename: DEFAULT_CSV_NAME,
		},
	}
}

func LoadTomlConfig() TomlDocument {
	var tomlDoc TomlDocument
	path := filepath.Join(DefaultPath(), DEFAULT_CONFIG_NAME)
	_, err := toml.DecodeFile(path, &tomlDoc)

	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("could not parse config.toml: %v", err)
		}

		// config file does not exist -> create config file
		createProjectDir()
		tomlDoc := DefaultTomlConfig()
		err := createTomlFile(tomlDoc)
		if err != nil {
			log.Fatal(err)
		}
	}

	return tomlDoc
}

func createProjectDir() {
	err := os.Mkdir(DefaultPath(), 0744)
	if err != nil {
		if os.IsExist(err) {
			log.Fatal("This should never happen. This Method should only be called if the project irectory does not exist yet and needs to be created!")
		}
		log.Fatal(err)
	}
}

func createTomlFile(config TomlDocument) error {
	path := filepath.Join(DefaultPath(), DEFAULT_CONFIG_NAME)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating config.toml: %v", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("error encoding toml config (%v): %v", config, err)
	}

	log.Printf("created default toml config (%s)\n", DEFAULT_CONFIG_NAME)
	return nil
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return filepath.Join(home, ".goalkeeper/")
}
