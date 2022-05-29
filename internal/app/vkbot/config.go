package vkbot

import "github.com/Kvertinum01/gomkpbot/internal/app/store"

type Config struct {
	Token      string `toml:"token"`
	GroupID    int    `toml:"group_id"`
	StrGroupID string `toml:"str_group_id"`
	Store      *store.Config
}
