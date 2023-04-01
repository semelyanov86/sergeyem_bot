package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
	"errors"
	"strconv"
	"strings"
)

const MsgLinksFromList = "Последние ссылки из списка 👉 "

type LinksFromListHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewLinksFromListHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) LinksFromListHandler {
	return LinksFromListHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h LinksFromListHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskListForLinks
}

func (h LinksFromListHandler) Handle(msg string, setting *settings.Setting) error {
	setting, err := h.settingsService.GetByUserName(setting.Username)
	if err != nil {
		return e.Wrap("failed to get settings", err)
	}
	if setting == nil {
		err = h.tg.SendMessage(setting.ChatId, NoSettingsError)
		if err != nil {
			return e.Wrap("error while sending response with no setting error", err)
		}
		return e.Wrap("ask for latestLinks with no settings", err)
	}
	err = h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		return e.Wrap("failed to change mode to root", err)
	}
	err = h.settingsService.SetContext(setting.Username, "")
	if err != nil {
		return e.Wrap("failed to clear context", err)
	}

	perPage, err := strconv.Atoi(setting.Context)
	if err != nil {
		perPage = 0
	}
	msgAttributes := strings.Split(msg, "|")

	linksError := links.ErrTokenNotExist
	listId, err := strconv.Atoi(msgAttributes[1])
	if err != nil {
		h.tg.SendMessage(setting.ChatId, "There was an error while converting data from context")
		return e.Wrap("converting context error", err)
	}
	latestLinks, err := h.linkService.GetLinksFromList(perPage, listId)
	if errors.Is(linksError, err) {
		h.tg.SendMessage(setting.ChatId, LinkAceTokenError)
		return e.Wrap("auth error with LinkAce", err)
	}

	if err != nil {
		return e.Wrap("error while getting latest latestLinks", err)
	}

	var text = "<i>" + MsgLinksFromList + msgAttributes[0] + "</i>"
	if len(latestLinks) < 0 {
		text = "Не найдено каких-либо ссылок 🤔 " + msgAttributes[0]
	}
	text = text + "\n"

	for i, link := range latestLinks {
		text = text + strconv.Itoa(i+1) + ". <b>" + link.Title + "</b>\n" + link.URL + "\n"
	}

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with linksFromList", err)
	}
	return nil
}
