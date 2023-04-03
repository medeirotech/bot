package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/liverday/medeiro-tech-bot/config"
)

var (
	gTenorUrl = "https://tenor.googleapis.com/v2/search"

	fridayTrigger  = "s e x t o u"
	fridayGifUrl   = "https://media.tenor.com/zGlEbV_bTnIAAAAC/kowalski-familia.gif"
	fallbackGifUrl = "https://media.tenor.com/RtJifRTjOHEAAAAC/dancing-random.gif"
)

var cfg config.Config

type GTenorMinimalReturn struct {
	Results []struct {
		ID                 string `json:"id"`
		Title              string `json:"title"`
		ContentDescription string `json:"content_description"`
		ContentRating      string `json:"content_rating"`
		H1Title            string `json:"h1_title"`
		Media              struct {
			Mp4 struct {
				Dims     []int   `json:"dims"`
				Preview  string  `json:"preview"`
				Size     int     `json:"size"`
				URL      string  `json:"url"`
				Duration float64 `json:"duration"`
			} `json:"mp4"`
			Gif struct {
				Size     int    `json:"size"`
				URL      string `json:"url"`
				Preview  string `json:"preview"`
				Dims     []int  `json:"dims"`
				Duration int    `json:"duration"`
			} `json:"gif"`
			Tinygif struct {
				Dims    []int  `json:"dims"`
				Size    int    `json:"size"`
				Preview string `json:"preview"`
				URL     string `json:"url"`
			} `json:"tinygif"`
		} `json:"media_formats"`
		BgColor    string        `json:"bg_color"`
		Created    float64       `json:"created"`
		Itemurl    string        `json:"itemurl"`
		URL        string        `json:"url"`
		Tags       []interface{} `json:"tags"`
		Flags      []interface{} `json:"flags"`
		Shares     int           `json:"shares"`
		Hasaudio   bool          `json:"hasaudio"`
		Hascaption bool          `json:"hascaption"`
		SourceID   string        `json:"source_id"`
		Composite  interface{}   `json:"composite"`
	} `json:"results"`
	Next string `json:"next"`
}

func FridayHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	cfg = config.GetConfig()

	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.Contains(strings.ToLower(m.Content), fridayTrigger) {
		return
	}

	message := &discordgo.MessageSend{
		Files: []*discordgo.File{},
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
		gifUrl = extractGifFromGTenor(gif, fallbackGifUrl)
	default:
		content = fmt.Sprintf("Calma família ainda não é sexta! Falta %d dia(s)", daysRemainingToFriday())

		gif = getRandomGif(time.Now().Weekday().String())
		gifUrl = extractGifFromGTenor(gif, fallbackGifUrl)
	}

	message.Content = content
	message.Files = append(message.Files, processGifUrl(gifUrl))

	_, err := s.ChannelMessageSendComplex(m.ChannelID, message)

	if err != nil {
		fmt.Println("Friday Handler - There was an exception when sending a message", err)
	}
}

func getRandomGif(search string) (result GTenorMinimalReturn) {
	req, err := http.NewRequest("GET", gTenorUrl, nil)
	if err != nil {
		fmt.Println("Cannot make a new http Request", err)
	}

	query := req.URL.Query()
	query.Add("random", "true")
	query.Add("q", search)
	query.Add("key", cfg.GTenorKey)
	query.Add("media_filter", "minimal")
	query.Add("limit", "1")

	req.URL.RawQuery = query.Encode()

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error on get a random gif", err)
	}

	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshall JSON", err)
	}

	return result
}

func extractGifFromGTenor(gTenor GTenorMinimalReturn, fallback string) string {
	if len(gTenor.Results) > 0 {
		return gTenor.Results[0].Media.Gif.URL
	}

	return fallback
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
