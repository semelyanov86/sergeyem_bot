package words

import "context"

type RepositoryInterface interface {
	GetRandomWords(ctx context.Context, perPage int) ([]Word, error)
	SaveWord(ctx context.Context, word *Word) error
	GetSettings(ctx context.Context) (WordSettings, error)
}
