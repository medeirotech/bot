package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/liverday/medeiro-tech-bot/config"
)

type ApiResponse struct {
	Data struct {
		Id        string `json:"id"`
		ShortLink string `json:"short_link"`
	} `json:"data"`
}

type ApiRequest struct {
	Link string `json:"link"`
}

var (
	res = make(chan *ApiResponse)
	err = make(chan error)
)

func worker(d string, cfg *config.Config) {
	url := "https://api.curto.io/v1/links"
	log.Printf("New Destination Received to shorten: %s\n", d)

	apiRequest := &ApiRequest{
		Link: strings.TrimSpace(d),
	}

	c := &http.Client{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(apiRequest)
	req, e := http.NewRequest(http.MethodPost, url, buf)

	if e != nil {
		err <- e
		return
	}

	req.Header.Set("X-Curto-Api-Key", cfg.CurtoApiKey)
	req.Header.Set("Content-Type", "application/json")

	r, e := c.Do(req)

	if e != nil {
		err <- e
		return
	}

	if e != nil {
		err <- e
		return
	}

	if r.StatusCode != 200 {
		err <- fmt.Errorf("the url creation failed: %d", r.StatusCode)
		return
	}

	defer r.Body.Close()
	body, e := io.ReadAll(r.Body)

	if e != nil {
		err <- e
		return
	}

	var apiResponse *ApiResponse

	e = json.Unmarshal(body, &apiResponse)

	if e != nil {
		err <- e
		return
	}

	res <- apiResponse
}

func ShortUrlHandler(cfg *config.Config) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		destination := i.ApplicationCommandData().Options[0].StringValue()

		go worker(destination, cfg)

		var content string
		select {
		case item := <-res:
			log.Printf("A url was created successfully, short link: %s\n", item.Data.ShortLink)
			content = fmt.Sprintf("🎉 Seu link curto é esse: %s", item.Data.ShortLink)
		case e := <-err:
			log.Printf("[ERROR] An error was thrown when creating a url: %s\n", e)
			content = fmt.Sprintf("Houve um erro: %s", e)
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: content,
			},
		})
	}
}
