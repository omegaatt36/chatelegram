package repository

import (
	"strings"
	"time"

	"github.com/omegaatt36/chatelegram/appmodule/telegram/usecase"
	"gopkg.in/telebot.v3"
)

// TelegramBot is implement of usecase.TelegramUseCase.
type TelegramBot struct {
	bot *telebot.Bot
}

var _ usecase.TelegramUseCase = &TelegramBot{}

// NewTelegramBot returns implement of usecase.TelegramUseCase.
func NewTelegramBot(bot *telebot.Bot) *TelegramBot {
	return &TelegramBot{bot: bot}
}

func ensureFormatting(text string) string {
	numDelimiters := strings.Count(text, "```")
	numSingleDelimiters := strings.Count(strings.Replace(text, "```", "", -1), "`")

	if (numDelimiters % 2) == 1 {
		text += "```"
	}
	if (numSingleDelimiters % 2) == 1 {
		text += "`"
	}

	return text
}

// SendAsLiveOutput sends message as live output.
func (b *TelegramBot) SendAsLiveOutput(chatID int64, feed <-chan string) error {
	var (
		sent     *telebot.Message
		lastResp string
		m        = message{feed: feed}
	)

	m.aggregate()

	send := func(register string) error {
		defer func() {
			lastResp = register
		}()

		if len(strings.Trim(register, "\n")) == 0 {
			time.Sleep(time.Microsecond * 50)
			return nil
		}

		if sent == nil {
			var err error

			sent, err = b.bot.Send(&telebot.Chat{ID: chatID}, register)
			return err
		}

		if register == lastResp {
			return nil
		}

		text := ensureFormatting(register)
		if _, err := b.bot.Edit(sent, text); err != nil {
			return err
		}

		return nil
	}

	for {
		select {
		case <-m.done:
			return send(m.str)
		default:
		}

		register := m.str

		if register == "" {
			continue
		}

		if err := send(register); err != nil {
			return err
		}
	}
}
