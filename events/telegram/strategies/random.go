package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
	"errors"
	"strconv"
)

const RandomCmd = "random"

const MsgRandom = "Список случайных слов 👇\n"

const AskForWordsToken = "Мы не нашли токен EasyWords или он не верный. Пожалуйста, предоставьте нам новый 😗"

type RandomHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordsService    words.WordService
}

func NewRandomHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) RandomHandler {
	return RandomHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordsService:    wordService,
	}
}

func (h RandomHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == RandomCmd
}

func (h RandomHandler) Handle(msg string, setting *settings.Setting) error {
	setting, err := h.settingsService.GetByUserName(setting.Username)
	if err != nil {
		return e.Wrap("failed to get settings", err)
	}
	if setting == nil {
		err = h.tg.SendMessage(setting.ChatId, NoSettingsError)
		if err != nil {
			return e.Wrap("error while sending response with no setting error", err)
		}
		return e.Wrap("ask for random words with no settings", err)
	}

	perPage, err := strconv.Atoi(h.meta.Message.CommandArguments())
	if err != nil {
		perPage = 0
	}

	var text = "<i>" + MsgRandom + "</i>"
	wordsError := words.ErrWordTokenNotExist
	latestWords, err := h.wordsService.GetRandomWords(perPage)
	if errors.Is(wordsError, err) {
		err := h.askWordsToken(setting)
		if err != nil {
			return e.Wrap("error while asking for easywords token", err)
		}
		return err
	}
	if len(latestWords) < 1 {
		text = "Не найдено каких-либо слов для изучения 🤔"
	}
	if err != nil {
		return e.Wrap("error while getting random words", err)
	}

	for i, word := range latestWords {
		text = text + strconv.Itoa(i+1) + ". <b>" + word.Original + "</b> - " + word.Translated + "\n"
	}

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with random words", err)
	}
	return nil
}

func (h RandomHandler) askWordsToken(setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskWordsToken)
	if err != nil {
		return e.Wrap("failed to change mode in AskWordsToken", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, AskForWordsToken); err != nil {
		return e.Wrap("error while sending response with words", err)
	}
	return nil
}
