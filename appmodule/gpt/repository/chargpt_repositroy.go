package repository

import (
	"context"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/omegaatt36/chatelegram/appmodule/gpt/usecase"
)

// GPT3Client is implement of usecase.GPTUseCase.
type GPT3Client struct {
	client gpt3.Client

	common
}

var _ usecase.GPTUseCase = &GPT3Client{}

// NewGPT3Client returns implement of usecase.GPTUseCase.
func NewGPT3Client(client gpt3.Client, options ...ClientOption) *GPT3Client {
	c := &GPT3Client{
		client: client,
		common: common{
			maxToken: 1000,
			engine:   gpt3.TextDavinci003Engine,
		},
	}

	for _, option := range options {
		option.injectOption(&c.common)
	}

	return c
}

// Stream asks GPT the question and receives answer.
func (c *GPT3Client) Stream(ctx context.Context, question string) (<-chan string, <-chan error) {
	res := make(chan string)
	errCh := make(chan error)
	go func() {
		defer close(res)
		defer close(errCh)
		err := c.client.CompletionStreamWithEngine(ctx,
			c.engine, gpt3.CompletionRequest{
				Prompt:      []string{question},
				MaxTokens:   gpt3.IntPtr(c.maxToken),
				Temperature: gpt3.Float32Ptr(0),
			}, func(resp *gpt3.CompletionResponse) {
				res <- resp.Choices[0].Text
			})
		if err != nil {
			errCh <- err
		}
	}()

	return res, errCh
}
