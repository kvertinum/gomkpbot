package duel

import "github.com/Kvertinum01/gomkpbot/internal/app/models"

type Member struct {
	Attacked bool
	Attack   int
	Protect  int
	Health   int
	IsWin    bool
	Model    *models.User
}

func NewMember(model *models.User) *Member {
	return &Member{
		Attacked: false,
		Attack:   0,
		Protect:  0,
		Health:   3,
		Model:    model,
	}
}

type Duel struct {
	Members       map[int]*Member
	NowWay        int
	AnotherMember int
	Ways          int
}
