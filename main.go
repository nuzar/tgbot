package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	// Token is telegram api token
	Token = os.Getenv("API_TOKEN")
	// URL is our service's webhook url
	URL = os.Getenv("WEBHOOK_URL")
)

// WebhookPath updates receiver router path
const WebhookPath = "/update"

func main() {
	if err := setWebHook(setWebHookReq{URL: URL + WebhookPath}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	r.POST(WebhookPath, handleNewUpdate)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func handleNewUpdate(c *gin.Context) {
	var update Update
	c.Bind(&update)

	log.Printf("received update: %#v", update)

	sendResponse(update.Message)
}

func sendResponse(m Message) {
	if err := reverseMessage(m); err != nil {
		log.Printf("reverseMessage error: %s", err)
	}
}

func reverseMessage(m Message) error {
	req := sendMessageReq{
		ChatID:           UnionIntString{int64: m.Chat.ID},
		Text:             Reverse(m.Text),
		ReplyToMessageID: m.MessageID,
	}
	return sendMessage(req)
}

// Reverse reverse string
func Reverse(in string) string {
	runes := []rune(in)
	low, high := 0, len(runes)-1
	for low < high {
		runes[low], runes[high] = runes[high], runes[low]
		low++
		high--
	}
	return string(runes)
}
