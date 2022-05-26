package vkbot

import (
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type Bot struct {
	api *vkapi.Api
}

func SetupBot(config *Config) error {
	api := vkapi.NewApi(config.Token)
	bot := &Bot{api: api}
	lp, err := vkapi.NewLongpoll(api, config.GroupID)
	if err != nil {
		return err
	}

	go lp.ListenNewMessages()

	for {
		select {
		case message := <-lp.LastMessage:
			err := bot.checkMessage(message)
			if err != nil {
				return err
			}
		case event := <-lp.LastEvent:
			log.Println(event)
		}
	}
}

func (bot *Bot) checkMessage(message vkapi.Message) error {
	switch message.Text {
	case "ping":
		return bot.api.Method("messages.send", map[string]interface{}{
			"user_id":   message.FromID,
			"random_id": 0,
			"message":   "pong",
		}, nil)
	}
	return nil
}
