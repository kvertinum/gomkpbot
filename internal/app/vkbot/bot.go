package vkbot

import (
	"fmt"
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/store"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type Bot struct {
	api     *vkapi.Api
	store   *store.Store
	config  *Config
	groupID int
	mention string
}

func SetupBot(config *Config) error {
	api := vkapi.NewApi(config.Token)
	bot := &Bot{
		api:     api,
		groupID: config.GroupID,
		config:  config,
		mention: fmt.Sprintf(
			"[club%v|@%s] ",
			config.GroupID, config.StrGroupID,
		),
	}
	lp, err := vkapi.NewLongpoll(api, config.GroupID)
	if err != nil {
		return err
	}

	if err := bot.configureStore(); err != nil {
		return err
	}

	go lp.ListenNewMessages()

	for {
		select {
		case message := <-lp.NewMessage:
			go bot.checkMessage(message)
		case event := <-lp.NewEvent:
			go bot.checkEvent(event)
		}
	}
}

func (bot *Bot) configureStore() error {
	st := store.New(bot.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	bot.store = st
	return nil
}

func (bot *Bot) checkMessage(message vkapi.Message) {
	if message.Action != nil {
		if message.Action.Type == "chat_invite_user" && message.Action.MemberID == -bot.groupID {
			if err := bot.send(message.PeerID, helloMessage); err != nil {
				log.Fatal(err)
			}
			return
		}
	}
	if message.PeerID >= 2000000000 {
		bot.checkChat(message)
	}
}

func (bot *Bot) checkEvent(event vkapi.LongpollMessage) {
	// In dev
}

func (bot *Bot) send(peerID int, message string) error {
	return bot.api.Method("messages.send", map[string]interface{}{
		"peer_id":   peerID,
		"random_id": 0,
		"message":   message,
	}, nil)
}
