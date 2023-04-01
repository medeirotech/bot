package handlers

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

var (
	fridayTrigger = "sextou"
	fridayGifUrl  = "https://media.tenor.com/zGlEbV_bTnIAAAAC/kowalski-familia.gif"
)

func FridayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.Contains(m.Content, fridayTrigger) {
		return
	}

	message := &discordgo.MessageSend{
		Files: []*discordgo.File{},
	}

	switch time.Now().Weekday() {
		case time.Friday:
			message.Content = "Sextouu família"
			message.Files = append(message.Files, processGifUrl(fridayGifUrl))
		default:
			message.Content = fmt.Sprintf("Calma família ainda não é sexta! Falta %d dia(s)", daysRemainingToFriday())
	}

	s.ChannelMessageSendComplex(m.ChannelID, message)
}

func processGifUrl(url string) *discordgo.File {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Bad request", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to get GIT", err)
	}

	gifFile := &discordgo.File{
		Name:   "sextou-familia.gif",
		Reader: bytes.NewReader(body),
	}

	return gifFile
}

func daysRemainingToFriday() int {
	today := time.Now()

	if today.Weekday() > time.Friday {
		return int(today.Weekday())
	} else {
		return int(time.Friday) - int(today.Weekday()) 
	}
}