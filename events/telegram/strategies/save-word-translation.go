package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
)

const MsgAskTranslation = "Введите перевод слова:"

type SaveWordTranslationHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordsService    words.WordService
}

func NewSaveWordTranslationsHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) SaveWordTranslationHandler {
	return SaveWordTranslationHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordsService:    wordService,
	}
}

func (h SaveWordTranslationHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskedWordOriginal
}

func (h SaveWordTranslationHandler) Handle(msg string, setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskedWordTranslation)
	if err != nil {
		return e.Wrap("failed to change mode to AskedWordOriginal", err)
	}

	err = h.settingsService.SetContext(setting.Username, msg)
	if err != nil {
		return e.Wrap("failed to change context", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, MsgAskTranslation); err != nil {
		return e.Wrap("error while sending message MsgAskTranslations", err)
	}
	return nil
}
