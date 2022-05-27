package vkbot

import (
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type Bot struct {
	api     *vkapi.Api
	groupID int
}

func SetupBot(config *Config) error {
	api := vkapi.NewApi(config.Token)
	bot := &Bot{api: api, groupID: -config.GroupID}
	lp, err := vkapi.NewLongpoll(api, config.GroupID)
	if err != nil {
		return err
	}

	go lp.ListenNewMessages()

	for {
		select {
		case message := <-lp.LastMessage:
			go bot.checkMessage(message)
		case event := <-lp.LastEvent:
			go bot.checkEvent(event)
		}
	}
}

func (bot *Bot) checkMessage(message vkapi.Message) {
	if message.Action != nil {
		if message.Action.Type == "chat_invite_user" && message.Action.MemberID == bot.groupID {
			if err := bot.send(message.PeerID, "здарова"); err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	switch message.Text {
	case "ping":
		if err := bot.send(message.PeerID, "pong"); err != nil {
			log.Fatal(err)
		}
	}
}

func (bot *Bot) checkEvent(event vkapi.LongpollMessage) {

}

func (bot *Bot) send(peerID int, message string) error {
	return bot.api.Method("messages.send", map[string]interface{}{
		"peer_id":   peerID,
		"random_id": 0,
		"message":   message,
	}, nil)
}
