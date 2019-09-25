package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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
	var update Update
	c.Bind(&update)

	log.Printf("received update: %#v", update)

	go sendResponse()
}

func sendResponse() {

}

type sendMessageReq struct {
	ChatID UnionIntString `json:"chat_id"`
}

// UnionIntString union type of int and string
type UnionIntString struct {
	int
	string
}

func (u UnionIntString) IsInt() bool {
	if u.int != 0 {
		return true
	}

	if u.string != "" {
		return false
	}

	return true
}

// MarshalJSON implement json.Marshaler interface
// marshal int first
func (u UnionIntString) MarshalJSON() ([]byte, error) {
	if !u.IsInt() {
		return []byte("\"" + u.string + "\""), nil

	}

	return []byte(strconv.Itoa(u.int)), nil
}

// UnmarshalJSON implement  json.Unmarshaler interface
func (u *UnionIntString) UnmarshalJSON(data []byte) error {
	switch data[0] {
	case '"':
		u.string = string(data[1 : len(data)-1])
		return nil
	default:
		if n, err := strconv.Atoi(string(data)); err != nil {
			return fmt.Errorf("%s is not UnionIntString", string(data))
		} else {
			u.int = n
		}
	}

	return nil
}

func sendMessage() {

}

func isHttpStatusOK(code int) bool {
	return code/200 == 1
}
