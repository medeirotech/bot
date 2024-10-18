package messages

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

	lastSlashIndex := strings.LastIndex(url, "/")
	filename := url[lastSlashIndex+1:]

	gifFile := &discordgo.File{
		Name:   filename,
		Reader: bytes.NewReader(body),
	}

	return gifFile
}

func daysRemainingTo(day time.Weekday) int {
	today := time.Now()

	remaining := 0

	for {
		if (int(today.Weekday()) % 7 == int(day)) {
			return remaining
		}

		remaining++
		today = today.Add(time.Hour * 24)
	}
}

func replyMessage(m *discordgo.MessageCreate, s *discordgo.Session, gifUrl string, content string) error {
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
	return err
}
