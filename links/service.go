package links

import (
	"bot/settings"
	"context"
	"errors"
	"time"
)

type LinkService struct {
	Repository RepositoryInterface
	Settings   *settings.Setting
	Config     settings.Config
}

var ErrTokenNotExist = errors.New("token not exist")

func (s LinkService) GetLatestLinks(perPage int) ([]Link, error) {
	var token = s.Settings.LinkaceToken
	var links []Link
	if token == "" {
		return links, ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if perPage < 1 {
		perPage = s.Config.LinksPerPage
	}
	return s.Repository.GetLatestLinks(ctx, perPage, 1)
}

func (s LinkService) GetLinksFromList(perPage int, listId int) ([]Link, error) {
	var token = s.Settings.LinkaceToken
	var links []Link
	if token == "" {
		return links, ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if perPage < 1 {
		perPage = s.Config.LinksPerPage
	}
	return s.Repository.GetLinksFromList(ctx, listId, perPage, 1)
}

func (s LinkService) GetAllLists() ([]List, error) {
	var token = s.Settings.LinkaceToken
	var lists []List
	if token == "" {
		return lists, ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return s.Repository.GetAllLists(ctx)
}

func (s LinkService) SaveLink(link *Link) error {
	var token = s.Settings.LinkaceToken
	if token == "" {
		return ErrTokenNotExist
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	return s.Repository.SaveLink(ctx, link)
}
