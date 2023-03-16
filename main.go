package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/omegaatt36/chatgpt-telegram/app"
	chatgpttelegram "github.com/omegaatt36/chatgpt-telegram/app/chatgpt-telegram"
	chatgpt "github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/repository"
	telegram "github.com/omegaatt36/chatgpt-telegram/appmodule/telegram/repository"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
	"gopkg.in/telebot.v3"
)

var config struct {
	telegramBotToken  string
	apiKey            string
	maxToken          int
	completionsEngine string
	timeout           int
	allowedUsers      []int64
}

// Main starts process in cli.
func Main(ctx context.Context) {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  config.telegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	cfg := openai.DefaultConfig(config.apiKey)
	cfg.HTTPClient = &http.Client{
		Timeout: time.Duration(time.Duration(config.timeout) * time.Second),
	}

	client := openai.NewClientWithConfig(cfg)

	service := chatgpttelegram.NewService(
		bot,
		telegram.NewTelegramBot(bot),
		chatgpt.NewOpenAIClient(client,
			chatgpt.WithMaxToken{MaxToken: config.maxToken},
			chatgpt.WithCompletionsEngine{Engine: config.completionsEngine},
		),
	)

	service.Start(ctx,
		chatgpttelegram.UseAllowedUsers{AllowedUsers: config.allowedUsers},
	)

	<-ctx.Done()
	log.Println("app stopping")
}

func main() {
	app := app.App{
		Main:  Main,
		Flags: []cli.Flag{},
	}

	app.Flags = append(app.Flags,
		&cli.StringFlag{
			Name:        "chatgpt-api-key",
			EnvVars:     []string{"CHATGPT_API_KEY"},
			Destination: &config.apiKey,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "telegram-bot-token",
			EnvVars:     []string{"TELEGRAM_BOT_TOKEN"},
			Destination: &config.telegramBotToken,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "chatgpt-max-token",
			EnvVars:     []string{"CHATGPT_MAX_TOKEN"},
			Destination: &config.maxToken,
			DefaultText: "3000",
			Value:       3000,
		},
		&cli.StringFlag{
			Name:        "chatgpt-completions-model",
			EnvVars:     []string{"CHATGPT_COMPLETIONS_MODEL"},
			Destination: &config.completionsEngine,
			DefaultText: "text-davinci-003",
			Value:       "text-davinci-003",
		},
		&cli.IntFlag{
			Name:        "chatgpt-timeout",
			EnvVars:     []string{"CHATGPT_TIMEOUT"},
			Destination: &config.timeout,
			DefaultText: "60",
			Value:       60,
		},
		&cli.MultiInt64Flag{
			Target: &cli.Int64SliceFlag{
				Name:    "allowed-users",
				EnvVars: []string{"ALLOWED_USERS"},
			},
			Value:       []int64{},
			Destination: &config.allowedUsers,
		},
	)

	app.Run()
}
