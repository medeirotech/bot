package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/liverday/medeiro-tech-bot/config"
)

var (
	sundayTrigger  = "domingo a noite"
	sundayGifUrl   = "https://media1.tenor.com/m/3udM407rgkQAAAAC/kikimogi-kiki.gif"
	sundayFallbackGifUrl = "https://media1.tenor.com/m/UAK1t55SoRUAAAAd/chaves-chavo.gif"
)


func SundayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	cfg = config.GetConfig()

	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.Contains(strings.ToLower(m.Content), sundayTrigger) {
		return
	}

	var content string
	var gif GTenorMinimalReturn
	var gifUrl string

	switch time.Now().Weekday() {
	case time.Sunday:
		content = "É hoje... Ele chegou..."
		gif = getRandomGif("o malvado")
		gifUrl = extractGifFromGTenor(gif, sundayGifUrl)
		
	case time.Saturday:
		content = "É amanhã... Ele nunca esteve tão perto..."
		gif = getRandomGif("prepare")
		gifUrl = extractGifFromGTenor(gif, sundayFallbackGifUrl)

	case time.Monday:
		content = "Ele foi embora ontem, mas se aproxima novamente..."
		gif = getRandomGif("monday mood")
		gifUrl = extractGifFromGTenor(gif, sundayFallbackGifUrl)

	default:
		content = fmt.Sprintf("Prepare-se, o evento canônico ocorrerá em %d dia(s)... Domingo a noite se aproxima.", daysRemainingToFriday())
		gif = getRandomGif("rezem")
		gifUrl = extractGifFromGTenor(gif, sundayFallbackGifUrl)
	}

	message := &discordgo.MessageSend{
		Content: content,
		Files:   []*discordgo.File{processGifUrl(gifUrl)},
		Reference: &discordgo.MessageReference{
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
		},
	}

	_, err := s.ChannelMessageSendComplex(m.ChannelID, message)

	if err != nil {
		fmt.Println("Sunday Handler - There was an exception when sending a message", err)
	}
}
