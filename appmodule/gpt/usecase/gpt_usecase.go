//go:generate mockgen -package=usecase -destination=gpt_usecase_mock.go . GPTUseCase

package usecase

import "context"

// GPTUseCase defines GPT send question use case.
type GPTUseCase interface {
	CompletionStream(ctx context.Context, question string) (<-chan string, <-chan error)
}
