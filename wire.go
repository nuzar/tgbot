//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
)

func InitBot() (*Bot, error) {
	wire.Build(NewBot, fromEnv, newTgBot, getUpdatesChan)
	return nil, nil
}
