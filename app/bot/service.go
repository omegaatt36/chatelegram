package chatelegram

import (
	"context"
	"errors"
	"log"
	"os"
	"unicode/utf8"

	gpt "github.com/omegaatt36/chatelegram/appmodule/gpt/usecase"
	telegram "github.com/omegaatt36/chatelegram/appmodule/telegram/usecase"
	"github.com/omegaatt36/chatelegram/src/health"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

// Service is GPT agent via telegram bot.
type Service struct {
	ctx context.Context

	allowedUsers []int64

	bot      *telebot.Bot
	telegram telegram.TelegramUseCase
	gpt      gpt.GPTUseCase
}

// NewService return Service with use cases.
func NewService(bot *telebot.Bot, tu telegram.TelegramUseCase, cu gpt.GPTUseCase) *Service {
	return &Service{
		bot:      bot,
		telegram: tu,
		gpt:      cu,
	}
}

func (s *Service) registerEndpoint() {
	s.bot.Handle("/start", s.handleStart)
	s.bot.Handle(telebot.OnText, s.handleTextCompletion)
}

func (s *Service) useMiddleware() {
	if len(s.allowedUsers) != 0 {
		s.bot.Use(middleware.Whitelist(s.allowedUsers...))
	}
}

// Start starts telegram bot service with context, and register stop event.
func (s *Service) Start(ctx context.Context, configs ...config) {
	s.ctx = ctx

	s.useMiddleware()
	s.registerEndpoint()

	go func() {
		log.Println("starting telegram bot")
		s.bot.Start()
	}()

	go func() {
		<-ctx.Done()
		log.Println("stopping telegram bot")
		s.bot.Stop()
		log.Println("telegram bot is stopped")
	}()

	go health.StartServer()
}

func (s Service) processTextCompeltion(chatID int64, question string) error {
	textCompletionStreamCh, errTextCompletionStreamChCh := s.gpt.Stream(s.ctx, question)

	errCh := make(chan error, 1)
	defer close(errCh)

	done := make(chan struct{}, 1)
	go func() {
		if err := s.telegram.SendAsLiveOutput(chatID, textCompletionStreamCh); err != nil {
			errCh <- err
		}
		done <- struct{}{}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-done:
			return nil
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			if err == nil {
				continue
			}

			return err
		case err, ok := <-errTextCompletionStreamChCh:
			if !ok {
				return nil
			}
			if err == nil {
				continue
			}

			return err
		}
	}
}

func (s *Service) handleTextCompletion(c telebot.Context) error {
	if utf8.RuneCountInString(c.Message().Text) == 0 {
		return nil
	}

	log.Printf("start(%d) user(%d) prompt(%s)\n",
		c.Message().ID, c.Message().Chat.ID, c.Message().Text)
	defer func() { log.Printf("done(%d)\n", c.Message().ID) }()

	err := s.processTextCompeltion(c.Message().Chat.ID, c.Message().Text)
	if err != nil {
		errMessage := err.Error()
		if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
			errMessage = "OpenAI is currently not responding, please try again later."
		}

		if ierr := c.Send(errMessage); ierr != nil {
			log.Fatalln(ierr)
		}
	}

	return err
}

func (s *Service) handleStart(c telebot.Context) error {
	return c.Send("wellcome to use GPT agent, please ask me something.")
}
