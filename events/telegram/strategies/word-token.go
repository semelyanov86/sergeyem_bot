package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
)

const MsgWordToken = `Токен успешно сохранён. Теперь вы можете дальше пользоваться сервисом EasyWords 👍`
const MsgErrorWordToken = "Невозможно сохранить токен. Произошла непредвиденная ошибка 😔 "

type WordTokenHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordService     words.WordService
}

func NewWordTokenHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) WordTokenHandler {
	return WordTokenHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordService:     wordService,
	}
}

func (h WordTokenHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskWordsToken
}

func (h WordTokenHandler) Handle(text string, setting *settings.Setting) error {
	setting.EasywordsToken = text
	err := h.settingsService.UpdateSetting(setting)
	if err != nil {
		err := h.tg.SendMessage(setting.ChatId, MsgErrorWordToken+"while saving setting")
		if err != nil {
			return err
		}
		return e.Wrap("failed to save new settings", err)
	}

	err = h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		return e.Wrap("failed to change mode to root", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgWordToken)
}
