package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Kvertinum01/gomkpbot/internal/app/store"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkbot"
)

var (
	configPath string
)

func init() {
	// Parse path to config from command args
	flag.StringVar(&configPath, "config-path", "configs/vkbot.toml", "path to config file")
}

func main() {
	flag.Parse()

	// Read .toml and write to config
	config := &vkbot.Config{Store: store.NewConfig()}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup bot
	if err := vkbot.SetupBot(config); err != nil {
		log.Fatal(err)
	}
}
