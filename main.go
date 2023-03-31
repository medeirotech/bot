package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/liverday/medeiro-tech-bot/config"
	"github.com/liverday/medeiro-tech-bot/handlers"
)

var (
	cfg      *config.Config
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "source",
			Description: "Get the source code of this bot",
		},
		{
			Name:        "short_url",
			Description: "Shorten an URL using curto.io engine",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "destination",
					Description: "Destination",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"source":    handlers.SourceHandler,
		"short_url": handlers.ShortUrlHandler,
	}
)

var bot *discordgo.Session

func contains(arr []string, target string) bool {
	for _, s := range arr {
		if s == target {
			return true
		}
	}

	return false
}

func isCorrectMessage(channelId string, messageId string) bool {
	return cfg.ChannelID == channelId && cfg.MessageID == messageId
}

func reactionAddHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	log.Printf("Adding reaction %s with the %s channelID and %s messageId\n", m.Emoji.Name, m.ChannelID, m.MessageID)
	if !isCorrectMessage(m.ChannelID, m.MessageID) {
		return
	}

	log.Printf("The message is correct, checking roles %+v\n", m.Member.Roles)

	if roleID, exists := cfg.RoleMap[m.Emoji.Name]; exists {
		if contains(m.Member.Roles, roleID) {
			return
		}

		log.Printf("Adding role %s to user %s of guild %s\n", roleID, m.UserID, m.GuildID)
		s.GuildMemberRoleAdd(m.GuildID, m.UserID, roleID)
	} else {
		log.Printf("A Role with the %s reaction was not found", m.Emoji.Name)
	}
}

func reactionRemoveHandler(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	log.Printf("New reaction %s removed with the %s channelID and %s messageId\n", m.Emoji.Name, m.ChannelID, m.MessageID)
	if !isCorrectMessage(m.ChannelID, m.MessageID) {
		return
	}

	if roleID, exists := cfg.RoleMap[m.Emoji.Name]; exists {
		log.Printf("Removing role %s from user %s of guild %s\n", roleID, m.UserID, m.GuildID)
		s.GuildMemberRoleRemove(m.GuildID, m.UserID, roleID)
	} else {
		log.Printf("A Role with the %s reaction was not found", m.Emoji.Name)
	}
}

func addHandlers() {
	bot.AddHandler(reactionAddHandler)
	bot.AddHandler(reactionRemoveHandler)
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
}

func addCommands() []*discordgo.ApplicationCommand {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		createdCommand, err := bot.ApplicationCommandCreate(bot.State.User.ID, cfg.GuildID, command)

		if err != nil {
			log.Fatalf("There was an error creating %v command: %v ", command.Name, err)
		}

		registeredCommands[i] = createdCommand
	}

	return registeredCommands
}

func removeCommands(commands []*discordgo.ApplicationCommand) {
	log.Println("Removing commands...")
	for _, command := range commands {
		err := bot.ApplicationCommandDelete(bot.State.User.ID, cfg.GuildID, command.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
		}
	}
}

func start() {
	log.Println("Starting MedeiroTech bot")
	var err error
	bot, err = discordgo.New(fmt.Sprintf("Bot %s", cfg.Token))

	if err != nil {
		log.Fatal("Error loading bot")
	}

	addHandlers()

	err = bot.Open()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	commands := addCommands()

	log.Println("The bot is running. Press CTRL-C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if cfg.RemoveCommands {
		removeCommands(commands)
	}
	log.Println("Gracefully shutting down.")
}

func main() {
	env := os.Getenv("ENVIRONMENT")

	if env != "PRODUCTION" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	config, err := config.Load()

	if err != nil {
		log.Fatal("Error loading config file")
	}

	cfg = config
	start()
}
