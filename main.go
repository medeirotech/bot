package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/liverday/medeiro-tech-reaction-bot/config"
)

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

func start() {
	log.Println("Starting MedeiroTech bot")
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.Token))

	if err != nil {
		log.Fatal("Error loading bot")
	}

	bot.AddHandler(reactionAddHandler)
	bot.AddHandler(reactionRemoveHandler)

	err = bot.Open()

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println("The bot is running")
}

var cfg *config.Config

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config, err := config.Load()

	if err != nil {
		log.Fatal("Error loading config file")
	}

	cfg = config
	start()

	<-make(chan struct{})
}
