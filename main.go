package main

import (
	"bintracker/internal/bintracker"
	"bintracker/internal/service"
	"fmt"
	"log"

	"github.com/pelletier/go-toml"
)

const credentialsFile = "credentials.json"
const configFile = "config.toml"

func main() {
	config, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var cfg bintracker.Config
	err = config.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Config", cfg)
	bins := service.GetBins(cfg)
	service.AddBinsToCalendar(cfg, bins)
	fmt.Println(bins)
}
