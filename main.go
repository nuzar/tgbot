package main

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/phuslu/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.DefaultLogger.Caller = 1

	bot, err := InitBot()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	var wg sync.WaitGroup

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		defer wg.Done()
		<-sigCh
		cancel()
		shutdown(bot.tgbot)
	}()

	bot.processUpdates(ctx)
	wg.Wait()
}
