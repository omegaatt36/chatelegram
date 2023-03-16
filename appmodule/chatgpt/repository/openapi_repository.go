package repository

import (
	"context"
	"fmt"
	"io"

	"github.com/omegaatt36/chatgpt-telegram/appmodule/chatgpt/usecase"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

// OpenAIClient is implement of usecase.ChatGPTUseCase.
type OpenAIClient struct {
	client *openai.Client

	common
}

var _ usecase.ChatGPTUseCase = &OpenAIClient{}

// NewOpenAIClient returns implement of usecase.ChatGPTUseCase.
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

// Stream asks ChatGPT the question and receives answer.
func (c *OpenAIClient) Stream(ctx context.Context, question string) (<-chan string, <-chan error) {
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
			fmt.Printf("CompletionStream error: %v\n", err)
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
