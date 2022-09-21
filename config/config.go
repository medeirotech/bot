package config

import (
	"os"
	"strconv"
)

type Config struct {
	Token          string
	GuildID        string
	ChannelID      string
	MessageID      string
	RemoveCommands bool
	RoleMap        map[string]string
}

var cfg *Config

func Load() (*Config, error) {
	cfg = new(Config)
	var err error

	cfg.Token = os.Getenv("TOKEN")
	cfg.GuildID = os.Getenv("GUILD_ID")
	cfg.ChannelID = os.Getenv("CHANNEL_ID")
	cfg.MessageID = os.Getenv("MESSAGE_ID")
	cfg.RemoveCommands, err = strconv.ParseBool(os.Getenv("REMOVE_COMMANDS"))

	if err != nil {
		return nil, err
	}

	cfg.RoleMap = map[string]string{
		"🪄": os.Getenv("FRONTEND_ROLE_ID"),
		"🏭": os.Getenv("BACKEND_ROLE_ID"),
		"🚀": os.Getenv("FULLSTACK_ROLE_ID"),
		"📱": os.Getenv("MOBILE_ROLE_ID"),
	}

	return cfg, nil
}
