package main

import (
	"os"
)

type Config struct {
	Token      string `json:"token"`
	WebHookURL string `json:"web_hook_url"`
	Port       string `json:"port"`
}

func fromEnv() Config {
	var cfg = Config{
		Token:      os.Getenv("TG_API_TOKEN"),
		WebHookURL: os.Getenv("TG_WEBHOOK_URL"),
		Port:       os.Getenv("PORT"),
	}
	return cfg
}
