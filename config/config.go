package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Token          string
	GuildID        string
	ChannelID      string
	MessageID      string
	GTenorKey      string
	RemoveCommands bool
	RoleMap        map[string]string
	CurtoApiKey    string
}

var cfg *Config

func Load() (*Config, error) {
	cfg = new(Config)
	var err error

	cfg.Token = os.Getenv("TOKEN")
	cfg.GuildID = os.Getenv("GUILD_ID")
	cfg.ChannelID = os.Getenv("CHANNEL_ID")
	cfg.MessageID = os.Getenv("MESSAGE_ID")
	cfg.GTenorKey = os.Getenv("GTENOR_KEY")
	cfg.RemoveCommands, err = strconv.ParseBool(os.Getenv("REMOVE_COMMANDS"))
	cfg.CurtoApiKey = os.Getenv("CURTO_API_KEY")

	if err != nil {
		return nil, err
	}

	cfg.RoleMap = map[string]string{
		"🪄": os.Getenv("FRONTEND_ROLE_ID"),
		"🏭": os.Getenv("BACKEND_ROLE_ID"),
		"🚀": os.Getenv("FULLSTACK_ROLE_ID"),
		"📱": os.Getenv("MOBILE_ROLE_ID"),
		"💣": os.Getenv("QA_ROLE_ID"),
	}

	return cfg, nil
}

func GetConfig() Config {
	if cfg == nil {
		log.Fatal("You must load the config")
	}

	return *cfg
}
