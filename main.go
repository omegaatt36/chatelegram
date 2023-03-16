package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/omegaatt36/chatelegram/app"
	chatelegram "github.com/omegaatt36/chatelegram/app/bot"
	gpt "github.com/omegaatt36/chatelegram/appmodule/gpt/repository"
	telegram "github.com/omegaatt36/chatelegram/appmodule/telegram/repository"
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

	service := chatelegram.NewService(
		bot,
		telegram.NewTelegramBot(bot),
		gpt.NewOpenAIClient(client,
			gpt.WithMaxToken{MaxToken: config.maxToken},
			gpt.WithCompletionsEngine{Engine: config.completionsEngine},
		),
	)

	service.Start(ctx,
		chatelegram.UseAllowedUsers{AllowedUsers: config.allowedUsers},
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
			Name:        "openai-api-key",
			EnvVars:     []string{"OPENAI_API_KEY"},
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
			Name:        "gpt-max-token",
			EnvVars:     []string{"GPT_MAX_TOKEN"},
			Destination: &config.maxToken,
			DefaultText: "3000",
			Value:       3000,
		},
		&cli.StringFlag{
			Name:        "gpt-completions-model",
			EnvVars:     []string{"GPT_COMPLETIONS_MODEL"},
			Destination: &config.completionsEngine,
			DefaultText: "text-davinci-003",
			Value:       "text-davinci-003",
		},
		&cli.IntFlag{
			Name:        "gpt-timeout",
			EnvVars:     []string{"GPT_TIMEOUT"},
			Destination: &config.timeout,
			DefaultText: "60",
			Value:       60,
		},
		&cli.MultiInt64Flag{
			Target: &cli.Int64SliceFlag{
				Name:    "telegram-allowed-users",
				EnvVars: []string{"TELEGRAM_ALLOWED_USERS"},
			},
			Value:       []int64{},
			Destination: &config.allowedUsers,
		},
	)

	app.Run()
}
