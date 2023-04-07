package words

import (
	"bot/settings"
	"context"
	"time"
)

type WordService struct {
	Repository RepositoryInterface
	Settings   *settings.Setting
	Config     settings.Config
}

func (s WordService) GetRandomWords(perPage int) ([]Word, error) {
	var token = s.Settings.EasywordsToken
	var words []Word
	if token == "" {
		return words, ErrWordTokenNotExist
	}
	if perPage < 1 {
		perPage = s.Config.WordsPerPage
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.Repository.GetRandomWords(ctx, perPage)
}

func (s WordService) SaveWord(word *Word) error {
	var token = s.Settings.EasywordsToken
	if token == "" {
		return ErrWordTokenNotExist
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.Repository.SaveWord(ctx, word)
}
