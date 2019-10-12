package main

import "fmt"

// BotAPI https://api.telegram.org/bot<token>/METHOD_NAME
const BotAPI = "https://api.telegram.org/bot%s/%s"

func getApiURI(method string) string {
	return fmt.Sprintf(BotAPI, Token, method)
}

// BotAPIResponse BotAPI response
// If ‘ok’ equals true, the request was successful and the result of the query can be found in the ‘result’ field.
// In case of an unsuccessful request, ‘ok’ equals false and the error is explained in the ‘description’.
// An Integer ‘error_code’ field is also returned, but its contents are subject to change in the future.
// Some errors may also have an optional field ‘parameters’ of the type ResponseParameters,
// which can help to automatically handle the error.
type BotAPIResponse struct {
	OK bool `json:"ok"`
	// TODO: what's result?
	Result             interface{}        `json:"result"`
	Description        string             `json:"description"`
	ErrorCode          int                `json:"error_code"`
	ResponseParameters ResponseParameters `json:"parameters"`
}

// ResponseParameters help to automatically handle the error
type ResponseParameters struct {
	// MigrateToChatID Optional. The group has been migrated to a supergroup with the specified identifier.
	// This number may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it is smaller than 52 bits,
	// so a signed 64 bit integer or double-precision float type are safe for storing this identifier.
	MigrateToChatID int64 `json:"migrate_to_chat_id"`
	//RetryAfter Optional. In case of exceeding flood control,
	// the number of seconds left to wait before the request can be repeated
	RetryAfter int `json:"retry_after"`
}

type Update struct {
	// Optional. New incoming message of any kind — text, photo, sticker, etc.
	Message Message `json:"message"`
}

// Message this object represents a message.
type Message struct {
	// Unique message identifier inside this chat
	MessageID int  `json:"message_id"`
	Chat      Chat `json:"chat"`
	// Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters.
	Text string `json:"text"`
	// Date the message was sent in Unix time
	Data int64 `json:"date"`
}

// Chat This object represents a chat.
type Chat struct {
	// Unique identifier for this chat. This number may be greater than 32 bits, but it is smaller than 52 bits
	ID    int64    `json:"id" validate:"required"`
	Type  ChatType `json:"chat" validate:"required"`
	Title string   `json:"title"`
}

// ChatType type of chat
type ChatType string

const (
	ChatTypePrivate    ChatType = "private"
	ChatTypeGroup      ChatType = "group"
	ChatTypeSuperGroup ChatType = "supergroup"
	ChatTypeChannel    ChatType = "channel"
)
