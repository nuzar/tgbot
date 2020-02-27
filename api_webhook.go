package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nuzar/tgbot/log"
)

type setWebHookReq struct {
	URL string `json:"url"`
	// TODO: tls settings
	Certificate    interface{} `json:"certificate,omitempty"`
	MaxConnections int         `json:"max_connections,omitempty"`
	// List the types of updates you want your bot to receive.
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

func setWebHook(req setWebHookReq) error {
	log.S.Infof("set up webhook: %v", req)

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

	if !isHTTPStatusOK(resp.StatusCode) {
		return fmt.Errorf("setWebHook failed [%d]: %s", resp.StatusCode, string(b))
	}

	var botAPIResp BotAPIResponse
	if err := json.Unmarshal(b, &botAPIResp); err != nil {
		return err
	}

	if !botAPIResp.OK {
		return fmt.Errorf("setWebHook failed %s", botAPIResp.Description)
	}

	return nil
}
