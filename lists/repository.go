package lists

import (
	"bot/lib/e"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type ListRepository struct {
	Url   string
	Token string
}

var ErrWrongUrl = errors.New("wrong easylist url")
var ErrWrongStatusCode = errors.New("wrong status code")
var ErrInvalidData = errors.New("invalid data in json")
var ErrTokenNotExist = errors.New("token not exist")
var ErrWrongList = errors.New("invalid listID")

type ErrorMessage struct {
	Title  string
	Detail string
}

type ErrorsResponse struct {
	Errors []ErrorMessage
}

func (l ListRepository) GetLists(ctx context.Context, perPage int) ([]List, error) {
	var listResponse Lists
	var lists []List
	if l.Token == "" {
		return lists, ErrTokenNotExist
	}
	if l.Url == "" {
		return lists, ErrWrongUrl
	}
	if perPage < 1 {
		perPage = 100
	}
	url := l.Url + "/api/v1/lists?page[number]=1&page[size]=" + strconv.Itoa(perPage)
	req, err := l.generateRequestFromUrl(ctx, url)
	if err != nil {
		return lists, e.Wrap("failed to generate request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return lists, ErrTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return lists, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return lists, e.Wrap("error while getting lists", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return lists, e.Wrap("error while reading body", ErrInvalidData)
	}
	err = json.Unmarshal(body, &listResponse)
	if err != nil {
		return lists, e.Wrap("error while parsing body", ErrInvalidData)
	}
	return listResponse.Data, nil
}

func (l ListRepository) GetItemsFromList(ctx context.Context, listId int, perPage int) ([]Item, error) {
	var itemsResponse Items
	var items []Item
	if l.Token == "" {
		return items, ErrTokenNotExist
	}
	if l.Url == "" {
		return items, ErrWrongUrl
	}
	if perPage < 1 {
		perPage = 100
	}

	url := l.Url + "/api/v1/lists/" + strconv.Itoa(listId) + "/items?page[number]=1&page[size]=" + strconv.Itoa(perPage)

	req, err := l.generateRequestFromUrl(ctx, url)
	if err != nil {
		return items, e.Wrap("failed to generate request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return items, ErrTokenNotExist
	}
	if res.StatusCode == http.StatusNotFound {
		return items, ErrWrongList
	}
	if res.StatusCode != http.StatusOK {
		return items, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return items, e.Wrap("error while getting items from list", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return items, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &itemsResponse)
	if err != nil {
		return items, e.Wrap("error while parsing body", err)
	}
	return itemsResponse.Data, nil
}

func (l ListRepository) generateRequestFromUrl(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Chat-Bot")
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer "+l.Token)
	return req, nil
}
