package vkbot

import (
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

func SetupBot(config *Config) error {
	api := vkapi.NewApi(config.Token)
	lp, err := vkapi.NewLongpoll(api, config.GroupID)
	if err != nil {
		return err
	}

	messages := make(chan vkapi.Message)
	go lp.ListenNewMessages(messages)

	for message := range messages {
		switch message.Text {
		case "ping":
			api.Method("messages.send", map[string]interface{}{
				"user_id":   message.FromID,
				"random_id": 0,
				"message":   "pong",
			}, nil)
		}
	}

	return nil
}
