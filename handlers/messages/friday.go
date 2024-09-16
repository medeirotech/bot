package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/liverday/medeiro-tech-bot/config"
)

var (
	fridayTrigger        = "s e x t o u"
	fridayGifUrl         = "https://media.tenor.com/zGlEbV_bTnIAAAAC/kowalski-familia.gif"
	fridayFallbackGifUrl = "https://media.tenor.com/RtJifRTjOHEAAAAC/dancing-random.gif"
)

func FridayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	cfg = config.GetConfig()
	if !strings.Contains(strings.ToLower(m.Content), fridayTrigger) {
		return
	}

	var content string
	var gif GTenorMinimalReturn
	var gifUrl string

	switch time.Now().Weekday() {
	case time.Friday:
		content = "Sextouu família"
		gif = getRandomGif(fridayTrigger)
		gifUrl = extractGifFromGTenor(gif, fridayGifUrl)

	case time.Thursday:
		content = "Quase, mas ainda não"
		gif = getRandomGif("quase-la")
		gifUrl = extractGifFromGTenor(gif, fridayFallbackGifUrl)

	default:
		content = fmt.Sprintf("Calma família ainda não é sexta! Falta %d dia(s)", daysReminingTo(time.Friday))
		gif = getRandomGif(time.Now().Weekday().String())
		gifUrl = extractGifFromGTenor(gif, fridayFallbackGifUrl)
	}

	err := replyMessage(m, s, gifUrl, content)

	if err != nil {
		fmt.Println("Friday Handler - There was an exception when sending a message", err)
	}
}
