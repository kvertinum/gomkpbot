package vkbot

import (
	"database/sql"
	"log"
	"strings"

	"github.com/Kvertinum01/gomkpbot/internal/app/models"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

const (
	helloMessage = "Здравствуйте! Для полной работы бота, назначьте его администратором беседы"
	helpMessage  = `
	/stat - Показывает вашу статистику
	/duel <user> - Вызывает человека на дуэль
	`
)

func (bot *Bot) checkChat(message vkapi.Message) {
	// Catch the bot's join into the chat
	if message.Action != nil {
		if message.Action.Type == "chat_invite_user" && message.Action.MemberID == -bot.groupID {
			if err := bot.send(message.PeerID, helloMessage); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	// Checking the existence of user in the database
	model, err := bot.store.User().FindByID(message.FromID)
	if err != nil {
		model, err = bot.checkDbErr(err, message.FromID, message.PeerID)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Parse message
	cmdArgs := strings.Split(message.Text, " ")
	cmdName := cmdArgs[0]
	var cmdValues []string
	if len(cmdArgs) > 1 {
		cmdValues = cmdArgs[1:]
	}

	// Setup route
	route := &Route{
		bot:       bot,
		message:   message,
		userModel: model,
		cmdValues: cmdValues,
	}

	// Check message
	switch cmdName {
	case "/help":
		route.helpCmd()
	case "/stat":
		route.statCmd()
	case "/duel":
		route.duelCmd()
	}
}

func (bot *Bot) checkDbErr(err error, userID int, peerID int) (*models.User, error) {
	if err == sql.ErrNoRows {
		model := &models.User{
			UserID:   userID,
			UserName: "no name",
			PeerID:   peerID,
			Wins:     0,
			Loses:    0,
		}
		return model, bot.store.User().Create(model)
	} else {
		return nil, err
	}
}
