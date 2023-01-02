package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func getChatGPTresponse(ctx context.Context, question string) string {
	c := gogpt.NewClient(os.Getenv("OPENAI_TOKEN"))

	maxtokens, err0 := strconv.Atoi(os.Getenv("OPENAI_MAXTOKENS"))

	if err0 != nil {
		fmt.Println("Error during conversion")
		return "MaxTokens Conversion Error happened!"
	}

	req := gogpt.CompletionRequest{
		Model:       "text-davinci-003",
		MaxTokens:   maxtokens,
		Prompt:      question,
		Temperature: 0,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return "You got an error!"
	} else {
		fmt.Println(resp.Choices[0].Text)

		return resp.Choices[0].Text
	}

}

func main() {
	ctx := context.Background()
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage: //message.Text refers to the text users typed in; ctx has to be passed down to the function
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(getChatGPTresponse(ctx, message.Text))).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
