package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/phuslu/log"
)

var (
	// BotToken is telegram api token
	BotToken = MustGetEnv("TG_API_TOKEN")
	// WebHookURL is our service's webhook url
	WebHookURL = os.Getenv("TG_WEBHOOK_URL")
	// PORT is listen port
	PORT = os.Getenv("PORT")
)

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, updatesCh, err := setup(ctx)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	log.Info().Msg("setup finish")
	defer shutdown(bot)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cancel()
	}()

	processUpdates(ctx, bot, updatesCh)
}

func setup(ctx context.Context) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	var err error

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, nil, err
	}
	log.Info().Msgf("authorized on account %s", bot.Self.UserName)

	var updates tgbotapi.UpdatesChannel
	if WebHookURL == "" {
		log.Info().Msg("setup long polling")
		updates, err = setUpPolling(ctx, bot)
	} else {
		log.Info().Msg("setup web hook")
		updates, err = setupWebHook(ctx, bot)
	}

	return bot, updates, err
}

func setUpPolling(_ context.Context, bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return updates, nil
}

func setupWebHook(_ context.Context, bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	log.Info().Msgf("webhook url: %s", WebHookURL)

	_, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	if info.LastErrorDate != 0 {
		log.Error().Msgf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/")
	go func() {
		port := "8080"
		if PORT != "" {
			port = PORT
		}
		log.Info().Msgf("listen on port %s", port)

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	return updates, nil
}

func processUpdates(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	handler := func(update tgbotapi.Update) {
		if update.Message == nil { // ignore any non-Message Updates
			log.Debug().Msg("received nil message")
			return
		}

		log.Info().
			Str("username", update.Message.From.UserName).
			Str("text", update.Message.Text).
			Msg("received update")

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reverse(update.Message.Text))
		if _, err := bot.Send(msg); err != nil {
			log.Error().Err(err).Msg("send message failed")
		}
	}

	for {
		select {
		case update := <-updates:
			handler(update)
		case <-ctx.Done():
			log.Info().Msg(ctx.Err().Error())
			for update := range updates {
				handler(update)
			}
			return
		}
	}
}

func shutdown(bot *tgbotapi.BotAPI) {
	log.Info().Msg("shutdown")
	if bot == nil {
		return
	}

	bot.StopReceivingUpdates()

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Error().Err(err).Msg("get web hook failed")
	}
	if info.IsSet() {
		log.Info().Msg("delete webhook")
		resp, err := bot.MakeRequest("deleteWebhook", nil)
		if err != nil {
			log.Error().Err(err).Msgf("delete webhook failed: %+v", resp)
		}
	}
}

// MustGetEnv get env or panic
func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Errorf("cannot get %s", key))
	}
	return val
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
