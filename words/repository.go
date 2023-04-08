package words

import (
	"bot/lib/e"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type WordRepository struct {
	Token string
	Url   string
}

var ErrWrongWordsUrl = errors.New("wrong Easywords url")
var ErrWrongWordsStatusCode = errors.New("wrong status code")
var ErrWordsInvalidData = errors.New("invalid data in json")
var ErrWordTokenNotExist = errors.New("token not exist")

type ErrorWordResponse struct {
	Message string
}

type WordResponse struct {
	Data []Word
}

func (w WordRepository) GetRandomWords(ctx context.Context, perPage int) ([]Word, error) {
	var wordResponse WordResponse
	if w.Token == "" {
		return []Word{}, ErrWordTokenNotExist
	}
	if w.Url == "" {
		return []Word{}, ErrWrongWordsUrl
	}
	req, err := w.generateRequestFromUrl(ctx, w.Url+"/api/random/"+strconv.Itoa(perPage), "", nil)
	if err != nil {
		return []Word{}, e.Wrap("failed to generate words request", err)
	}
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return []Word{}, ErrWordTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return []Word{}, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongWordsStatusCode)
	}
	if err != nil {
		return []Word{}, e.Wrap("Error while getting words", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Word{}, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &wordResponse)
	if err != nil {
		return []Word{}, e.Wrap("error while parsing body", err)
	}
	return wordResponse.Data, nil
}

func (w WordRepository) GetSettings(ctx context.Context) (WordSettings, error) {
	var settingsResponse WordSettingsResponse
	if w.Token == "" {
		return WordSettings{}, ErrWordTokenNotExist
	}

	if w.Url == "" {
		return WordSettings{}, ErrWrongWordsUrl
	}

	req, err := w.generateRequestFromUrl(ctx, w.Url+"/api/settings/", "", nil)
	if err != nil {
		return WordSettings{}, e.Wrap("failed to generate settings request", err)
	}
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode == http.StatusUnauthorized {
		return WordSettings{}, ErrWordTokenNotExist
	}
	if res.StatusCode != http.StatusOK {
		return WordSettings{}, e.Wrap(fmt.Sprintf("Wrong status code: %d", res.StatusCode), ErrWrongWordsStatusCode)
	}
	if err != nil {
		return WordSettings{}, e.Wrap("Error while getting settings", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return WordSettings{}, e.Wrap("error while reading body", err)
	}
	err = json.Unmarshal(body, &settingsResponse)
	if err != nil {
		return WordSettings{}, e.Wrap("error while parsing body", err)
	}
	return settingsResponse.Data, nil
}

func (w WordRepository) SaveWord(ctx context.Context, word *Word) error {
	var errResp ErrorWordResponse
	wordBody, err := json.Marshal(word)
	if err != nil {
		return e.Wrap("can not marshal model to json", err)
	}

	req, err := w.generateRequestFromUrl(ctx, w.Url+"/api/words", http.MethodPost, bytes.NewBuffer(wordBody))
	if err != nil {
		return e.Wrap("can not create request", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusUnauthorized {
		return ErrWordTokenNotExist
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
		return e.Wrap(errResp.Message, ErrWordsInvalidData)
	}
	if res.StatusCode != http.StatusCreated {
		return e.Wrap(fmt.Sprintf("wrong status code %d", res.StatusCode), ErrWrongWordsStatusCode)
	}
	return nil
}

func (w WordRepository) generateRequestFromUrl(ctx context.Context, url string, method string, body io.Reader) (*http.Request, error) {
	if method == "" {
		method = http.MethodGet
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Chat-Bot")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+w.Token)
	return req, nil
}
