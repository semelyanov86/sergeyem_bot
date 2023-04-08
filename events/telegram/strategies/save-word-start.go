package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
)

const SaveWordStartCmd = "saveword"

const MsgSaveStart = "Введите оригинальное значение слова:"

type SaveWordStartHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordsService    words.WordService
}

func NewSaveWordStartHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) SaveWordStartHandler {
	return SaveWordStartHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordsService:    wordService,
	}
}

func (h SaveWordStartHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == SaveWordStartCmd
}

func (h SaveWordStartHandler) Handle(msg string, setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskedWordOriginal)
	if setting.EasywordsToken == "" {
		return h.askWordsToken(setting)
	}
	if err != nil {
		return e.Wrap("failed to change mode to AskedWordOriginal", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, MsgSaveStart); err != nil {
		return e.Wrap("error while sending message MsgSaveStart", err)
	}
	return nil
}

func (h SaveWordStartHandler) askWordsToken(setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskWordsToken)
	if err != nil {
		return e.Wrap("failed to change mode in AskWordsToken", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, AskForWordsToken); err != nil {
		return e.Wrap("error while sending response with words", err)
	}
	return nil
}
