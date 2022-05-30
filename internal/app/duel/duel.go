package duel

type Member struct {
	Attack  int
	Protect int
	IsWin   bool
}

type Duel struct {
	Members       map[int]*Member
	NowWay        int
	AnotherMember int
	Ways          int
}
