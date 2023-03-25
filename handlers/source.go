package handlers

import "github.com/bwmarrin/discordgo"

func SourceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "O link para o repositório é esse: https://github.com/liverday/medeiro-tech-bot",
		},
	})
}
