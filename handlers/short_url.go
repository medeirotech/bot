package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
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

func worker(d string) {
	url := "https://api.curto.io/v1/urls"
	log.Printf("New Destination Received to shorten: %s\n", d)

	apiRequest := &ApiRequest{
		Link: d,
	}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(apiRequest)
	r, e := http.Post(url, "application/json", buf)

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

func ShortUrlHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	destination := i.ApplicationCommandData().Options[0].StringValue()

	go worker(destination)

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
