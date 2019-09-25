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
