package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nuzar/tgbot/log"
	"go.uber.org/zap/zapcore"
)

var (
	// BotToken is telegram api token
	BotToken = MustGetEnv("TG_API_TOKEN")
	// WebHookURL is our service's webhook url
	WebHookURL = os.Getenv("TG_WEBHOOK_URL")
	PORT       = os.Getenv("PORT")
)

func main() {
	var err error

	log.Init(zapcore.DebugLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot, updatesCh, err := setup(ctx)
	if err != nil {
		log.L.Fatal(err)
	}
	log.L.Info("setup finish")
	defer shutdown(bot)

	go processUpdates(ctx, bot, updatesCh)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	cancel()
}

func setup(ctx context.Context) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	var err error

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		return nil, nil, err
	}
	log.L.Infof("Authorized on account %s", bot.Self.UserName)

	var updates tgbotapi.UpdatesChannel
	if WebHookURL == "" {
		log.L.Info("setup long polling")
		updates = setUpPolling(ctx, bot)
	} else {
		log.L.Info("setup web hook")
		updates = setupWebHook(ctx, bot)
	}

	return bot, updates, nil
}

func setUpPolling(_ context.Context, bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.L.Fatal(err)
	}
	return updates
}

func setupWebHook(_ context.Context, bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	log.L.Info("webhook url: ", WebHookURL)

	_, err := bot.SetWebhook(tgbotapi.NewWebhook(WebHookURL))
	if err != nil {
		log.L.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.L.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.L.Errorf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/")
	go func() {
		port := "8080"
		if PORT != "" {
			_, err := strconv.Atoi(PORT)
			if err != nil {
				log.L.Error("invalid port ", PORT)
			} else {
				port = PORT
			}
		}
		log.L.Info("listen on port ", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.L.Fatal(err)
		}
	}()

	return updates
}

func processUpdates(ctx context.Context, bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}
			log.L.
				With("username", update.Message.From.UserName).
				With("text", update.Message.Text).
				Info("received update")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reverse(update.Message.Text))
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = bot.Send(msg)
		case <-ctx.Done():
			log.L.Info(ctx.Err())
			return
		}
	}
}

func shutdown(bot *tgbotapi.BotAPI) {
	log.L.Info("shutdown")
	if bot == nil {
		return
	}

	bot.StopReceivingUpdates()

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.L.Error(err)
	}
	if info.IsSet() {
		log.L.Info("delete webhook")
		resp, err := bot.MakeRequest("deleteWebhook", url.Values{})
		if err != nil {
			log.L.With("resp", resp, "err", err).Error("delete webhook failed")
		}
	}

	_ = log.L.Sync()
}

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
