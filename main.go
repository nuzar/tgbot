package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type setWebHookReq struct {
	URL string `json:"url"`
	// TODO: tls settings
	Certificate    interface{} `json:"certificate,omitempty"`
	MaxConnections int         `json:"max_connections,omitempty"`
	// List the types of updates you want your bot to receive.
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

func setWebHook(req setWebHookReq) error {
	log.Printf("set up webhook: %v", req)

	const method = "setWebhook"
	uri := getApiURI(method)

	b, _ := json.Marshal(req)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !isHttpStatusOK(resp.StatusCode) {
		return fmt.Errorf("setWebHook failed [%d]: %s", resp.StatusCode, string(b))
	}

	var botApiResp BotAPIResponse
	if err := json.Unmarshal(b, &botApiResp); err != nil {
		return err
	}

	if !botApiResp.OK {
		return fmt.Errorf("setWebHook failed %s", botApiResp.Description)
	}

	return nil
}

func handleNewUpdate(c *gin.Context) {
	var req map[string]interface{}
	c.Bind(&req)
	log.Printf("received updates: %s", req)
}

func isHttpStatusOK(code int) bool {
	return code/200 == 1
}
