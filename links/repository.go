package links

import (
	"bot/lib/e"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type LinkRepository struct {
	Url   string
	Token string
}

var ErrWrongUrl = errors.New("wrong linkace url")
var ErrWrongStatusCode = errors.New("wrong status code")
var ErrInvalidData = errors.New("invalid data in json")

type LinkAceResponse[T Link | List] struct {
	Total int
	To    int
	Data  []T
}

type ErrorResponse struct {
	Message string
}

func (l2 LinkRepository) SaveLink(ctx context.Context, l *Link) error {
	var errResp ErrorResponse
	linkBody, err := json.Marshal(l)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l2.Url+"/api/v1/links", bytes.NewBuffer(linkBody))
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Chat-Bot")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l2.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return ErrTokenNotExist
	}
	if res.StatusCode == http.StatusUnprocessableEntity {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &errResp)
		if err != nil {
			return err
		}
		return e.Wrap(errResp.Message, ErrInvalidData)
	}

	if res.StatusCode != http.StatusOK {
		return e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return e.Wrap("error while storing links", err)
	}
	defer res.Body.Close()

	return nil
}

func (l2 LinkRepository) GetLatestLinks(ctx context.Context, limit int, page int) ([]Link, error) {
	var linksResponse LinkAceResponse[Link]
	if l2.Token == "" {
		return []Link{}, ErrTokenNotExist
	}
	if l2.Url == "" {
		return []Link{}, ErrWrongUrl
	}
	if page < 1 {
		page = 1
	}

	urlStr, err := l2.getUrlFromParams("/api/v1/links", limit, page)
	if err != nil {
		return []Link{}, e.Wrap("failed to generate url", err)
	}

	req, err := l2.generateRequestFromUrl(ctx, urlStr)
	if err != nil {
		return []Link{}, e.Wrap("failed to generate request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return []Link{}, ErrTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return []Link{}, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return []Link{}, e.Wrap("error while getting links", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Link{}, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &linksResponse)
	if err != nil {
		return []Link{}, e.Wrap("error while parsing body", err)
	}
	return linksResponse.Data, nil
}

func (l2 LinkRepository) GetLinksFromList(ctx context.Context, listId int, limit int, page int) ([]Link, error) {
	var linksResponse LinkAceResponse[Link]
	if l2.Token == "" {
		return []Link{}, ErrTokenNotExist
	}
	if l2.Url == "" {
		return []Link{}, ErrWrongUrl
	}
	if page < 1 {
		page = 1
	}

	urlStr, err := l2.getUrlFromParams("/api/v1/lists/"+strconv.Itoa(listId)+"/links", limit, page)
	if err != nil {
		return []Link{}, e.Wrap("failed to generate url", err)
	}

	req, err := l2.generateRequestFromUrl(ctx, urlStr)
	if err != nil {
		return []Link{}, e.Wrap("failed to generate request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return []Link{}, ErrTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return []Link{}, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return []Link{}, e.Wrap("error while getting links", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Link{}, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &linksResponse)
	if err != nil {
		return []Link{}, e.Wrap("error while parsing body", err)
	}
	return linksResponse.Data, nil
}

func (l2 LinkRepository) GetAllLists(ctx context.Context) ([]List, error) {
	var linksResponse LinkAceResponse[List]
	if l2.Token == "" {
		return []List{}, ErrTokenNotExist
	}
	if l2.Url == "" {
		return []List{}, ErrWrongUrl
	}

	urlStr, err := l2.getUrlFromParams("/api/v1/lists", -1, 1)
	if err != nil {
		return []List{}, e.Wrap("failed to generate list url", err)
	}

	req, err := l2.generateRequestFromUrl(ctx, urlStr)
	if err != nil {
		return []List{}, e.Wrap("failed to generate request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return []List{}, ErrTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return []List{}, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongStatusCode)
	}
	if err != nil {
		return []List{}, e.Wrap("error while getting lists", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []List{}, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &linksResponse)
	if err != nil {
		return []List{}, e.Wrap("error while parsing body", err)
	}
	return linksResponse.Data, nil
}

func (l2 LinkRepository) getUrlFromParams(resource string, perPage int, page int) (string, error) {
	baseURL := l2.Url
	params := url.Values{}
	params.Add("per_page", strconv.Itoa(perPage))
	params.Add("page", strconv.Itoa(page))

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", ErrWrongUrl
	}
	u.Path = resource
	u.RawQuery = params.Encode()
	return fmt.Sprintf("%v", u), nil
}

func (l2 LinkRepository) generateRequestFromUrl(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Chat-Bot")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l2.Token)
	return req, nil
}

func NewLinkRepository(url string, token string) LinkRepository {
	return LinkRepository{Url: url, Token: token}
}
