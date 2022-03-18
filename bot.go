package main

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/phuslu/log"
)

type Bot struct {
	tgbot   *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewBot(tgbot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) *Bot {
	return &Bot{
		updates: updates,
		tgbot:   tgbot,
	}
}

func (b *Bot) processUpdates(ctx context.Context) {
	for {
		select {
		case update := <-b.updates:
			if update.Message == nil || update.Message.Text == "" {
				log.Debug().Msg("received non-text message")
				continue
			}

			log.Info().
				Str("username", update.Message.From.UserName).
				Str("text", update.Message.Text).
				Msg("received update")

			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				reverse(update.Message.Text),
			)
			if _, err := b.tgbot.Send(msg); err != nil {
				log.Error().Err(err).Msg("send message failed")
			}
		case <-ctx.Done():
			return
		}
	}
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
