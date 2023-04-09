package lists

import (
	"bot/settings"
	"context"
	"time"
)

type ListService struct {
	Repository RepositoryInterface
	Settings   *settings.Setting
	Config     settings.Config
}

func (s ListService) GetAllLists(perPage int) ([]List, error) {
	var token = s.Settings.EasylistToken
	var lists []List
	if token == "" {
		return lists, ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if perPage < 1 {
		perPage = s.Config.ListsPerPage
	}
	return s.Repository.GetLists(ctx, perPage)
}

func (s ListService) GetItemsFromList(listId int) ([]Item, error) {
	var token = s.Settings.EasylistToken
	var items []Item
	if token == "" {
		return items, ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	perPage := s.Config.ListsPerPage

	return s.Repository.GetItemsFromList(ctx, listId, perPage)
}
