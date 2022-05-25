package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkbot"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/vkbot.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := &vkbot.Config{}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := vkbot.SetupBot(config); err != nil {
		log.Fatal(err)
	}
}
