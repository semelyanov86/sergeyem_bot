package words

import (
	"context"
)

type MockRepository struct {
	Token string
}

func (w MockRepository) GetRandomWords(ctx context.Context, perPage int) ([]Word, error) {
	words := []Word{
		{
			Id:         1,
			Original:   "Hamburg",
			Translated: "Гамбург",
			Views:      3,
			Language:   "DE",
			Starred:    false,
		},
		{
			Id:         2,
			Original:   "Stabil",
			Translated: "Стабильно",
			Views:      4,
			Language:   "DE",
			Starred:    false,
		},
	}
	return words, nil
}

func (w MockRepository) SaveWord(ctx context.Context, word *Word) error {
	return nil
}
