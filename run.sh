#!/bin/bash
BIN=./tgbot

export API_TOKEN="123:123"
export WEBHOOK_URL="https://example.com"

go build &&
    ${BIN}
