package vkbot

import (
	"fmt"
	"log"

	"github.com/Kvertinum01/gomkpbot/internal/app/duel"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

type PayloadRoute struct {
	bot     *Bot
	message vkapi.Message
	payload map[string]string
}

func (p *PayloadRoute) checkPayload() {
	// Check payload
	cmd, ok := p.payload["cmd"]
	if !ok {
		if err := p.bot.send(
			p.message.PeerID, "Используйте только кнопки!",
		); err != nil {
			log.Fatal(err)
		}
	}
	switch cmd {
	case "accept_duel":
		p.acceptDuel()
	}
}

func (p *PayloadRoute) acceptDuel() {
	// Accept duel
	userID, ok := p.bot.waitDuel[p.message.FromID]
	if !ok {
		// Check user id
		if err := p.bot.send(
			p.message.PeerID, "Вы не являетесь соперником",
		); err != nil {
			log.Fatal(err)
		}
		return
	}
	// Start duel
	delete(p.bot.waitDuel, p.message.FromID)
	duelID := len(p.bot.duels) + 1

	// Get models
	firstModel, err := p.bot.store.User().FindByID(
		p.message.PeerID, p.message.FromID,
	)
	if err != nil {
		log.Fatal(err)
	}
	secondModel, err := p.bot.store.User().FindByID(
		p.message.PeerID, p.message.FromID,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Add duel to bot.duels
	p.bot.duels[duelID] = &duel.Duel{
		Members: map[int]*duel.Member{
			userID:           duel.NewMember(firstModel),
			p.message.FromID: duel.NewMember(secondModel),
		},
		NowWay:        userID,
		AnotherMember: p.message.FromID,
		Ways:          0,
	}

	// Create keyboard
	kjson, err := createKeyboard(duelID, "attack")
	if err != nil {
		log.Fatal(err)
	}

	// Start duel
	answer := fmt.Sprintf(
		"Дуэль начинается. Атакует: [id%v|%s]",
		userID, firstModel.UserName,
	)
	if err := p.bot.api.Method("messages.send", map[string]interface{}{
		"peer_id":   p.message.PeerID,
		"random_id": 0,
		"message":   answer,
		"keyboard":  kjson,
	}, nil); err != nil {
		log.Fatal(err)
	}
}

func createKeyboard(duelID int, wayType string) (string, error) {
	// Creaate keyboard to attack
	parts := map[int]string{
		1: "Голова",
		2: "Живот",
	}

	k := vkapi.NewKeyboard(false, true)
	for i := 1; i <= 2; i++ {
		k.Add(vkapi.NewCallbackButton(
			parts[i], fmt.Sprintf(
				"{\"way\": \"%v\", \"type\": \"%s\", \"duel_id\": \"%v\"}",
				i, wayType, duelID,
			), "negative",
		))
	}
	k.NewLine()
	return k.GetJson()
}
