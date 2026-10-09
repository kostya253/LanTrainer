package main

import (
	"log/slog"
	"os"

	botserver "github.com/kostya253/LanTrainer/internal/bot-server"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
	slog.Info("Loaded .env file")

	if os.Getenv("BOT_TOKEN") == "" || os.Getenv("BOT_NAME") == "" || os.Getenv("AUTHORIZED_USER") == "" {
		slog.Error("Missing required environment variables")
		os.Exit(1)
	}
}

func main() {
	slog.Info("Starting bot server")
	botServer := botserver.NewBotServer()
	botServer.Start()
}
