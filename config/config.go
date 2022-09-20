package config

import (
	"os"
)

type Config struct {
	Token     string
	ChannelID string
	MessageID string
	RoleMap   map[string]string
}

var cfg *Config

func Load() (*Config, error) {
	cfg = new(Config)

	cfg.Token = os.Getenv("TOKEN")
	cfg.ChannelID = os.Getenv("CHANNEL_ID")
	cfg.MessageID = os.Getenv("MESSAGE_ID")

	cfg.RoleMap = map[string]string{
		"🪄": os.Getenv("FRONTEND_ROLE_ID"),
		"🏭": os.Getenv("BACKEND_ROLE_ID"),
		"🚀": os.Getenv("FULLSTACK_ROLE_ID"),
		"📱": os.Getenv("MOBILE_ROLE_ID"),
	}

	return cfg, nil
}
