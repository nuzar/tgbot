package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nuzar/tgbot/log"
)

var (
	// BotToken is telegram api token
	BotToken = MustGetEnv("TG_API_TOKEN")
	// WebHookURL is our service's webhook url
	WebHookURL = MustGetEnv("TG_WEBHOOK_URL")
)

// WebhookPath updates receiver router path
const WebhookPath = "/update"

func main() {
	if err := setWebHook(setWebHookReq{URL: WebHookURL + WebhookPath}); err != nil {
		log.S.Fatal(err)
	}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "ok")
	})
	r.POST(WebhookPath, handleNewUpdate)

	log.S.Fatal(r.Run()) // listen and serve on 0.0.0.0:8080
}

func handleNewUpdate(c *gin.Context) {
	var update Update
	err := c.Bind(&update)
	if err != nil {
		log.S.Errorf("invalid update: %s", err)
	}

	log.S.Infof("received update: %#v", update)

	sendResponse(update.Message)
}

func sendResponse(m Message) {
	if err := sendReverseMsg(m); err != nil {
		log.S.Errorf("reverseMessage error: %s", err)
	}
}

func sendReverseMsg(m Message) error {
	req := sendMessageReq{
		ChatID:           UnionIntString{int64: m.Chat.ID},
		Text:             reverse(m.Text),
		ReplyToMessageID: m.MessageID,
	}
	return sendMessage(req)
}

// Reverse reverse string
func reverse(in string) string {
	runes := []rune(in)
	low, high := 0, len(runes)-1
	for low < high {
		runes[low], runes[high] = runes[high], runes[low]
		low++
		high--
	}
	return string(runes)
}

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Errorf("cannot get %s", key))
	}
	return val
}
