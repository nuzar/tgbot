package main

import (
	"fmt"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/phuslu/log"
)

func newTgBot(cfg Config) (*tgbotapi.BotAPI, error) {
	hc := &http.Client{
		Timeout: time.Second * 3,
	}

	bot, err := tgbotapi.NewBotAPIWithClient(cfg.Token, tgbotapi.APIEndpoint, hc)
	if err != nil {
		return nil, fmt.Errorf("create tg bot failed: %w", err)
	}
	log.Info().Msgf("authorized on account %s", bot.Self.UserName)

	if err := setupWebHook(bot, cfg); err != nil {
		return nil, err
	}

	return bot, nil
}

func setupWebHook(bot *tgbotapi.BotAPI, cfg Config) error {
	log.Info().Msgf("webhook url: %s", cfg.WebHookURL)

	wh, _ := tgbotapi.NewWebhook(cfg.WebHookURL)
	if _, err := bot.Request(wh); err != nil {
		return fmt.Errorf("request with webhook config failed: %w", err)
	}

	return nil
}

func getUpdatesChan(cfg Config, bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	updates := bot.ListenForWebhook("/")

	port := "8080"
	if cfg.Port != "" {
		port = cfg.Port
	}
	log.Info().Msgf("listen on port %s", port)

	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	return updates, nil
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
		resp, err := bot.MakeRequest("deleteWebhook", nil)
		if err != nil {
			log.Error().Err(err).Msgf("delete webhook failed: %+v", resp)
		}
		log.Info().Msg("webhook deleted")
	}
}
