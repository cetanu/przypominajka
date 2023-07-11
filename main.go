package main

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(0)
	var configPath string
	pflag.StringVar(&configPath, "config", "config.yaml", "config file")
	pflag.Parse()
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	var config any
	if err := viper.Unmarshal(&config, nil); err != nil {
		log.Fatalln(err)
	}
	log.Println(config)
}
