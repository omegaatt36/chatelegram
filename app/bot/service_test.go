package chatelegram

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	gpt "github.com/omegaatt36/chatelegram/appmodule/gpt/usecase"
	telegram "github.com/omegaatt36/chatelegram/appmodule/telegram/usecase"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	s := assert.New(t)

	var (
		ctx            = context.TODO()
		question       = "test"
		chatID   int64 = 1

		ch            = make(chan string)
		readOnlyCh    = func(c chan string) <-chan string { return c }(ch)
		readOnlyErrCh = make(<-chan error)
	)

	go func() {
		for _, r := range "answer" {
			ch <- string(r)
		}
	}()

	controller := gomock.NewController(t)
	mockTelegram := telegram.NewMockTelegramUseCase(controller)
	{
		mockTelegram.EXPECT().SendAsLiveOutput(chatID, readOnlyCh).
			Times(1).Return(nil)
	}

	mockGPT := gpt.NewMockGPTUseCase(controller)
	{
		mockGPT.EXPECT().Stream(ctx, question).Times(1).Return(readOnlyCh, readOnlyErrCh)
	}
	service := NewService(nil, mockTelegram, mockGPT)
	service.ctx = ctx
	s.NoError(service.processTextCompletion(chatID, question))
}
