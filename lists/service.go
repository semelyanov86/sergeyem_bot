package lists

import (
	"bot/settings"
	"context"
	"strconv"
	"time"
)

const MsgItems = "Товары из выбранного вами списка:"

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

func (s ListService) GenerateMessageFromItems(items []Item) string {
	if len(items) < 1 {
		return "Не найдено каких-либо элементов 🤔"
	}
	var text = "<i>" + MsgItems + "</i>\n"
	var postFix = " "
	for key, item := range items {
		if item.Attributes.Quantity > 0 {
			postFix = " (" + strconv.Itoa(item.Attributes.Quantity) + " " + item.Attributes.QuantityType + ")"
		} else {
			postFix = " "
		}
		text = text + strconv.Itoa(key+1) + ". <b>" + item.Attributes.Name + "</b>" + postFix + "\n" + item.Attributes.Description + "\n"
	}
	return text
}
