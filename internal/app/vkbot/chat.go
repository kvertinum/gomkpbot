package vkbot

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/models"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

func (bot *Bot) checkChat(message vkapi.Message) {
	model, err := bot.store.User().FindByID(message.FromID)
	if err != nil {
		model, err = bot.checkDbErr(err, message.FromID, message.PeerID)
		if err != nil {
			log.Fatal(err)
		}
	}
	switch message.Text {
	case bot.mention + "stat":
		if err := bot.send(message.PeerID, fmt.Sprintf(
			"Имя: %s\nПобеды: %v\nПоражения: %v",
			model.UserName, model.Wins, model.Loses,
		)); err != nil {
			log.Fatal(err)
		}
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
