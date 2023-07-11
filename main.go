package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Hi []string `yaml:"hi"`
}

func main() {
	log.SetFlags(0)
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()
	log.Println(readConfig(configPath))
}

func readConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatalln(err)
	}

	return &config, nil
}
