package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/url"
)

const MsgSaveLink = "Пожалуйста, выберите список, к которому мы прикрепим ссылку 👇\n"

type LinkSaveHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewLinkSaveHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) LinkSaveHandler {
	return LinkSaveHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h LinkSaveHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	if mode != Root {
		return false
	}
	return isCreateLinkCmd(h.meta.Message.Text)
}

func (h LinkSaveHandler) Handle(msg string, setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, AskList)
	if err != nil {
		h.tg.SendMessage(setting.ChatId, ModeChangeFailed)
		return e.Wrap("failed to change mode to AskList", err)
	}
	err = h.settingsService.SetContext(setting.Username, msg)
	if err != nil {
		_ = h.tg.SendMessage(setting.ChatId, ContextChangeFailed)
		return e.Wrap("failed to change context", err)
	}
	lists, err := h.linkService.GetAllLists()
	if len(lists) < 1 {
		h.tg.SendMessage(setting.ChatId, MsgNoList)
		return e.Wrap("No lists. Please create one first", err)
	}

	allRows := GenerateListsButtons(lists, false)
	message := tgbotapi.NewMessage(setting.ChatId, MsgSaveLink)

	keyboard := tgbotapi.NewReplyKeyboard(
		allRows...,
	)
	message.ReplyMarkup = keyboard
	if err := h.tg.Send(message); err != nil {
		return e.Wrap("Failed to send a keyboard", err)
	}
	return nil
}

func isCreateLinkCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
