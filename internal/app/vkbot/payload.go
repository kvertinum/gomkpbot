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
	userID, ok := p.bot.waitDuel[p.message.FromID]
	if !ok {
		if err := p.bot.send(
			p.message.PeerID, "Вы не являетесь соперником",
		); err != nil {
			log.Fatal(err)
		}
		return
	}
	delete(p.bot.waitDuel, p.message.FromID)
	duelID := len(p.bot.duels) + 1

	firstModel, err := p.bot.store.User().FindByID(userID)
	if err != nil {
		log.Fatal(err)
	}
	secondModel, err := p.bot.store.User().FindByID(p.message.FromID)
	if err != nil {
		log.Fatal(err)
	}

	p.bot.duels[duelID] = &duel.Duel{
		Members: map[int]*duel.Member{
			userID:           duel.NewMember(firstModel),
			p.message.FromID: duel.NewMember(secondModel),
		},
		NowWay:        userID,
		AnotherMember: p.message.FromID,
		Ways:          0,
	}
	kjson, err := createAttackKeyboard(duelID)
	if err != nil {
		log.Fatal(err)
	}

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

func createAttackKeyboard(duelID int) (string, error) {
	parts := map[int]string{
		1: "Голова",
		2: "Живот",
		3: "Руки",
		4: "Ноги",
	}

	k := vkapi.NewKeyboard(false, true)
	for i := 1; i <= 4; i++ {
		k.Add(vkapi.NewCallbackButton(
			parts[i], fmt.Sprintf(
				"{\"way\": \"%v\", \"type\": \"attack\", \"duel_id\": \"%v\"}",
				i, duelID,
			), "negative",
		))
		if i == 2 {
			k.NewLine()
		}
	}
	k.NewLine()
	return k.GetJson()
}
