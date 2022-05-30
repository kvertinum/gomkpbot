package vkbot

import (
	"github.com/Kvertinum01/gomkpbot/internal/app/duel"
	"github.com/Kvertinum01/gomkpbot/internal/app/store"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type Bot struct {
	api      *vkapi.Api
	store    *store.Store
	config   *Config
	groupID  int
	duels    map[int]*duel.Duel
	waitDuel map[int]int
}

func SetupBot(config *Config) error {
	// Init Api struct
	api := vkapi.NewApi(config.Token)
	// Init Bot struct
	bot := &Bot{
		api:      api,
		groupID:  config.GroupID,
		config:   config,
		duels:    make(map[int]*duel.Duel),
		waitDuel: make(map[int]int),
	}
	// Init Longpoll struct
	lp, err := vkapi.NewLongpoll(api, config.GroupID)
	if err != nil {
		return err
	}

	// Configure store
	if err := bot.configureStore(); err != nil {
		return err
	}

	// Start listening messages
	go lp.ListenNewMessages()

	// Start infinite loop
	for {
		select {
		// Waiting message or event
		case message := <-lp.NewMessage:
			go bot.checkMessage(message)
		case event := <-lp.NewEvent:
			go bot.checkEvent(event)
		}
	}
}

func (bot *Bot) configureStore() error {
	// Create and check store
	st := store.New(bot.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	bot.store = st
	return nil
}

func (bot *Bot) checkMessage(message vkapi.Message) {
	if message.PeerID >= 2000000000 {
		bot.checkChat(message)
	}
}

func (bot *Bot) send(peerID int, message string) error {
	return bot.api.Method("messages.send", map[string]interface{}{
		"peer_id":   peerID,
		"random_id": 0,
		"message":   message,
	}, nil)
}
