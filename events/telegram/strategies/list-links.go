package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const ListLinksCmd = "listlinks"

const MsgListLinks = "Выберите список, в котором вы хотите посмотреть ссылки ⌨️\n"

const MsgNoList = "Нет списков. Пожалуйста, создайте его в сервисе LinkAce. 😏"

type ListLinksHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewListLinksHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) ListLinksHandler {
	return ListLinksHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h ListLinksHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == ListLinksCmd
}

func (h ListLinksHandler) Handle(msg string, setting *settings.Setting) error {
	setting, err := h.settingsService.GetByUserName(setting.Username)
	if err != nil {
		return e.Wrap("failed to get settings", err)
	}
	if setting == nil {
		err = h.tg.SendMessage(setting.ChatId, NoSettingsError)
		if err != nil {
			return e.Wrap("error while sending response with no setting error", err)
		}
		return e.Wrap("ask for listLinks with no settings", err)
	}
	err = h.settingsService.ChangeMode(setting.Username, AskListForLinks)
	if err != nil {
		h.tg.SendMessage(setting.ChatId, ModeChangeFailed)
		return e.Wrap("failed to change mode", err)
	}

	perPage, err := strconv.Atoi(h.meta.Message.CommandArguments())
	if err != nil {
		perPage = 0
	}
	err = h.settingsService.SetContext(setting.Username, string(perPage))
	if err != nil {
		h.tg.SendMessage(setting.ChatId, ContextChangeFailed)
		return e.Wrap("failed to save context", err)
	}

	lists, err := h.linkService.GetAllLists()
	linksError := links.ErrTokenNotExist

	if errors.Is(linksError, err) {
		h.askLinksToken(setting)
	}
	if len(lists) < 1 {
		h.tg.SendMessage(setting.ChatId, MsgNoList)
		return e.Wrap("No lists. Please create one first", err)
	}

	allRows := GenerateListsButtons(lists, true)
	message := tgbotapi.NewMessage(setting.ChatId, MsgListLinks)

	keyboard := tgbotapi.NewReplyKeyboard(
		allRows...,
	)
	message.ReplyMarkup = keyboard
	if err := h.tg.Send(message); err != nil {
		return e.Wrap("Failed to send a keyboard", err)
	}

	return nil
}

func (h ListLinksHandler) askLinksToken(setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, LinksToken)
	if err != nil {
		return e.Wrap("failed to change mode in askLinksToken", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, AskForLinksToken); err != nil {
		return e.Wrap("error while sending response with links", err)
	}
	return nil
}
