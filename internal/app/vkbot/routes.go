package vkbot

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/Kvertinum01/gomkpbot/internal/app/models"
	"github.com/Kvertinum01/gomkpbot/internal/app/vkapi"
)

const (
	waitDuelMsg = "Вы вызвали соперника на дуэль"
)

type Route struct {
	bot       *Bot
	message   vkapi.Message
	userModel *models.User
	cmdValues []string
}

func (r *Route) helpCmd() {
	// Answer to "/help"
	if err := r.bot.send(r.message.PeerID, helpMessage); err != nil {
		log.Fatal(err)
	}
}

func (r *Route) statCmd() {
	// Answer to "/stat"
	if err := r.bot.send(r.message.PeerID, fmt.Sprintf(
		"Имя: %s\nПобеды: %v\nПоражения: %v",
		r.userModel.UserName, r.userModel.Wins, r.userModel.Loses,
	)); err != nil {
		log.Fatal(err)
	}
}

func (r *Route) duelCmd() {
	// Answer to "/duel"
	if r.cmdValues != nil {
		// Create regex
		re := regexp.MustCompile(`\[id(\d+)\|.+\]`)
		res := re.FindStringSubmatch(r.cmdValues[0])
		if len(res) < 2 {
			r.sendNeedArgs()
			return
		}

		// Convert to int
		strUserID := res[1]
		userID, err := strconv.Atoi(strUserID)
		if err != nil {
			r.sendNeedArgs()
			return
		}

		// Check user id
		if userID == r.message.FromID {
			if err := r.bot.send(
				r.message.PeerID, "Вы не можете вызать на дуэль себя",
			); err != nil {
				log.Fatal(err)
			}
			return
		}

		// Add to waiting
		r.bot.waitDuel[userID] = r.message.FromID

		// Start timer
		go func(userID int, waitUserID, peerID int) {
			timer := time.NewTimer(time.Minute)
			<-timer.C

			_, ok := r.bot.waitDuel[waitUserID]
			if !ok {
				return
			}
			delete(r.bot.waitDuel, waitUserID)
			userModel, err := r.bot.store.User().FindByID(
				peerID, userID,
			)
			if err != nil {
				log.Fatal(err)
			}
			waitUserModel, err := r.bot.store.User().FindByID(
				peerID, waitUserID,
			)
			if err != nil {
				log.Fatal(err)
			}
			message := fmt.Sprintf(
				"[id%v|%s], пользователь [id%v|%s] не принял дуэль в течении минуты и она отменяется",
				userID, userModel.UserName,
				waitUserID, waitUserModel.UserName,
			)
			if err := r.bot.send(peerID, message); err != nil {
				log.Fatal(err)
			}
		}(
			r.message.FromID,
			userID,
			r.message.PeerID,
		)

		// Create keyboard object
		k := vkapi.NewKeyboard(false, true)
		k.Add(vkapi.NewTextButton(
			"Принять",
			"{\"cmd\": \"accept_duel\"}",
			"negative",
		))
		k.NewLine()
		kjson, err := k.GetJson()
		if err != nil {
			log.Fatal(err)
		}

		// Send answer
		if err := r.bot.api.Method("messages.send", map[string]interface{}{
			"peer_id":   r.message.PeerID,
			"random_id": 0,
			"message":   waitDuelMsg,
			"keyboard":  kjson,
		}, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		r.sendNeedArgs()
	}
}

func (r *Route) nameCmd() {
	// Change user name
	if r.cmdValues != nil {
		// Check value lenght
		if len(r.cmdValues[0]) > 15 {
			if err := r.bot.send(
				r.message.PeerID, "Имя должно быть не больше 15 симвволов",
			); err != nil {
				log.Fatal(err)
			}
			return
		}

		// Update in database
		if err := r.bot.store.User().NameByID(
			r.message.PeerID, r.message.FromID, r.cmdValues[0],
		); err != nil {
			log.Fatal(err)
		}

		answer := fmt.Sprintf("Вы сменили своё имя на %s", r.cmdValues[0])
		if err := r.bot.send(
			r.message.PeerID, answer,
		); err != nil {
			log.Fatal(err)
		}
	} else {
		r.sendNeedArgs()
	}
}

func (r *Route) sendNeedArgs() {
	// Answer when using the command incorrectly
	if err := r.bot.send(
		r.message.PeerID, "Эта команда требует аргументов",
	); err != nil {
		log.Fatal(err)
	}
}
