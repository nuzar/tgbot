package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nuzar/tgbot/log"
)

type sendMessageReq struct {
	// Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	ChatID UnionIntString `json:"chat_id" validate:"required"`
	// Text of the message to be sent
	Text string `json:"text" validate:"required"`
	// Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
	ParseMode string `json:"parse_mode"`
	// Disables link previews for links in this message
	DisableWebPagePreview bool `json:"disable_web_page_preview"`
	// Sends the message silently. Users will receive a notification with no sound.
	DisableNotification bool `json:"disable_notification"`
	// If the message is a reply, ID of the original message
	ReplyToMessageID int `json:"reply_to_message_id"`
	// Additional interface options. A JSON-serialized object for an inline keyboard, custom reply keyboard, instructions to remove reply keyboard or to force a reply from the user.
	// ReplyMarkup string `json:"reply_markup"`
}

func sendMessage(req sendMessageReq) error {
	const method = "sendMessage"
	log.S.Debugf("sendMessage: %v", req)

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
		return fmt.Errorf("sendMessage failed [%d]: %s", resp.StatusCode, string(b))
	}

	var botAPIResp BotAPIResponse
	if err := json.Unmarshal(b, &botAPIResp); err != nil {
		return err
	}

	if !botAPIResp.OK {
		return fmt.Errorf("sendMessage failed %s", botAPIResp.Description)
	}

	return nil
}
