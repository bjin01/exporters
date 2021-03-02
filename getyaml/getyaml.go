package getyaml

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		ApiUrl   string `yaml:"apiurl"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Ports    string `yaml:"port"`
	} `yaml:"server"`

	Method1 struct {
		Methodname   string                 `yaml:"methodname"`
		InputVars    map[string]interface{} `yaml:"input_map"`
		Outvariables []string               `yaml:"out_variablenames"`
		Filters      map[string]interface{} `yaml:"filters"`
		Outputformat struct {
			Byvariable bool `yaml:"by_variable"`
			Byset      bool `yaml:"by_dataset"`
			Json       bool `yaml:"json"`
		} `yaml:"output_format"`
	} `yaml:"method1"`

	Method2 struct {
		Methodname   string                 `yaml:"methodname"`
		InputVars    map[string]interface{} `yaml:"input_map"`
		Outvariables []string               `yaml:"out_variablenames"`
		Filters      map[string]interface{} `yaml:"filters"`
		Outputformat struct {
			Byvariable bool `yaml:"by_variable"`
			Byset      bool `yaml:"by_dataset"`
			Json       bool `yaml:"json"`
		} `yaml:"output_format"`
	} `yaml:"method2"`

	Finalmethod struct {
		Methodname string `yaml:"methodname"`
		Options    struct {
			Meth2dependmeth1 bool `yaml:"meth2_depend_meth1"`
		} `yaml:"options"`

		InputVars    map[string]interface{} `yaml:"input_map"`
		Outvariables []string               `yaml:"out_variablenames"`
		Filters      map[string]interface{} `yaml:"filters"`
		Outputformat struct {
			Byvariable bool `yaml:"by_variable"`
			Byset      bool `yaml:"by_dataset"`
			Json       bool `yaml:"json"`
		} `yaml:"output_format"`
	} `yaml:"finalmethod"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
