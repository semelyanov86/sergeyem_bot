package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"bot/words"
	"strings"
)

const MsgStoreWord = "Слово успешно сохранено. Не забывайте учить его 😎"

const MsgStoreWordError = "Произошла ошибка при сохранении слова! 🤯"

type StoreWordHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	wordsService    words.WordService
}

func NewStoreWordHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, wordService words.WordService) StoreWordHandler {
	return StoreWordHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		wordsService:    wordService,
	}
}

func (h StoreWordHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskedWordLanguage
}

func (h StoreWordHandler) Handle(msg string, setting *settings.Setting) error {
	currentContext := setting.Context
	contextData := strings.Split(currentContext, "|")
	wordModel := words.Word{
		Id:         0,
		Original:   contextData[0],
		Translated: contextData[1],
		Views:      0,
		Language:   msg,
		Starred:    false,
	}
	err := h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		_ = h.tg.SendMessage(setting.ChatId, ModeChangeFailed)
		return e.Wrap("failed to change mode to Root", err)
	}
	err = h.settingsService.SetContext(setting.Username, "")
	if err != nil {
		h.tg.SendMessage(setting.ChatId, ContextChangeFailed)
		return e.Wrap("failed to change context", err)
	}
	err = h.wordsService.SaveWord(&wordModel)
	if err != nil {
		h.tg.SendMessage(setting.ChatId, MsgStoreWordError+" "+err.Error())
		return e.Wrap("failed to save a word", err)
	}
	return h.tg.SendMessage(setting.ChatId, MsgStoreWord)
}
