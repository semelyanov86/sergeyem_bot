package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const MsgAskWordLanguage = "Выберите язык оригинала:"

type SaveWordLanguageHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordsService    words.WordService
}

func NewSaveWordLanguageHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) SaveWordLanguageHandler {
	return SaveWordLanguageHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordsService:    wordService,
	}
}

func (h SaveWordLanguageHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskedWordTranslation
}

func (h SaveWordLanguageHandler) Handle(msg string, setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskedWordLanguage)
	if err != nil {
		return e.Wrap("failed to change mode to AskedWordLanguage", err)
	}
	curContext := setting.Context
	err = h.settingsService.SetContext(setting.Username, curContext+"|"+msg)
	if err != nil {
		return e.Wrap("failed to change context", err)
	}
	wordSettings, err := h.wordsService.GetSettings()
	if err != nil {
		h.tg.SendMessage(setting.ChatId, err.Error())
		return e.Wrap("can not get settings from EasyWords", err)
	}
	message := tgbotapi.NewMessage(setting.ChatId, MsgAskWordLanguage)
	allRows := h.generateLanguageButtons(wordSettings.LanguagesList)

	keyboard := tgbotapi.NewReplyKeyboard(
		allRows...,
	)
	message.ReplyMarkup = keyboard

	if err := h.tg.Send(message); err != nil {
		return e.Wrap("error while sending message MsgAskWordLanguage", err)
	}
	return nil
}

func (h SaveWordLanguageHandler) generateLanguageButtons(languages []string) [][]tgbotapi.KeyboardButton {
	row := make([]tgbotapi.KeyboardButton, 0)
	allRows := make([][]tgbotapi.KeyboardButton, 0)

	for i, language := range languages {
		if i > 0 && i%3 == 0 {
			rows := tgbotapi.NewKeyboardButtonRow(row...)
			allRows = append(allRows, rows)
			row = make([]tgbotapi.KeyboardButton, 0)
		}
		title := language

		row = append(row, tgbotapi.NewKeyboardButton(title))
	}
	if len(row) > 0 {
		rows := tgbotapi.NewKeyboardButtonRow(row...)
		allRows = append(allRows, rows)
	}
	return allRows
}
