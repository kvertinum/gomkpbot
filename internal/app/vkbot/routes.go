package vkbot

import (
	"fmt"
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/models"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type Route struct {
	bot       *Bot
	message   vkapi.Message
	userModel *models.User
	cmdValues []string
}

func (r *Route) helpCmd() {
	if err := r.bot.send(r.message.PeerID, helpMessage); err != nil {
		log.Fatal(err)
	}
}

func (r *Route) statCmd() {
	if err := r.bot.send(r.message.PeerID, fmt.Sprintf(
		"Имя: %s\nПобеды: %v\nПоражения: %v",
		r.userModel.UserName, r.userModel.Wins, r.userModel.Loses,
	)); err != nil {
		log.Fatal(err)
	}
}

func (r *Route) duelCmd() {
	if r.cmdValues != nil {

	} else {
		if err := r.bot.send(
			r.message.PeerID, "Эта команда требуег аргументов",
		); err != nil {
			log.Fatal(err)
		}
	}
}
