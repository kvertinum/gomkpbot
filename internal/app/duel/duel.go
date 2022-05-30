package duel

import "github.com/Kvertinum01/gomkpbot/internal/app/models"

type Member struct {
	Attack  int
	Protect int
	IsWin   bool
	Model   *models.User
}

type Duel struct {
	Members       map[int]*Member
	NowWay        int
	AnotherMember int
	Ways          int
}
