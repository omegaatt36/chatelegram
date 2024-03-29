package repository

import (
	"context"
	"io"

	"github.com/omegaatt36/chatelegram/appmodule/gpt/usecase"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

// OpenAIClient is implement of usecase.GPTUseCase.
type OpenAIClient struct {
	client *openai.Client

	common
}

var _ usecase.GPTUseCase = &OpenAIClient{}

// NewOpenAIClient returns implement of usecase.GPTUseCase.
func NewOpenAIClient(client *openai.Client, options ...ClientOption) *OpenAIClient {
	c := &OpenAIClient{
		client: client,
		common: common{
			maxToken: 1000,
			engine:   openai.GPT3TextDavinci003,
		},
	}

	for _, option := range options {
		option.injectOption(&c.common)
	}

	return c
}

// CompletionStream asks GPT the question and receives answer.
func (c *OpenAIClient) CompletionStream(ctx context.Context, question string) (<-chan string, <-chan error) {
	res := make(chan string)
	errCh := make(chan error)
	req := openai.CompletionRequest{
		Model:     c.engine,
		MaxTokens: c.maxToken,
		Prompt:    question,
		Stream:    true,
	}

	go func() {
		defer close(res)
		defer close(errCh)
		stream, err := c.client.CreateCompletionStream(ctx, req)
		if err != nil {
			errCh <- errors.Wrapf(err, "CompletionStream error")
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				errCh <- errors.Wrapf(err, "Stream error")
				return
			}

			res <- response.Choices[0].Text
		}
	}()

	return res, errCh
}
