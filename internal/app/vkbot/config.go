package vkbot

type Config struct {
	Token   string `toml:"token"`
	GroupID int    `toml:"group_id"`
}
