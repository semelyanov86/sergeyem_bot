package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
	"errors"
	"strconv"
)

const LinksCmd = "links"

const MsgLinks = "Ваши последние сохранённые ссылки из всех категорий 👇\n"

const AskForLinksToken = "Мы не нашли токен LinkAce или он не верный. Пожалуйста, предоставьте нам новый 😗"

type LinksHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
	linkService     links.LinkService
}

func NewLinksHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client, linkService links.LinkService) LinksHandler {
	return LinksHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h LinksHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == LinksCmd
}

func (h LinksHandler) Handle(msg string, setting *settings.Setting) error {
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

	perPage, err := strconv.Atoi(h.meta.Message.CommandArguments())
	if err != nil {
		perPage = 0
	}

	var text = "<i>" + MsgLinks + "</i>"
	linksError := links.ErrTokenNotExist
	latestLinks, err := h.linkService.GetLatestLinks(perPage)
	if errors.Is(linksError, err) {
		err := h.askLinksToken(setting)
		if err != nil {
			return e.Wrap("error while asking for linkace token", err)
		}
	}
	if len(latestLinks) < 1 {
		text = "Не найдено каких-либо ссылок 🤔"
	}
	if err != nil {
		return e.Wrap("error while getting latest latestLinks", err)
	}

	for i, link := range latestLinks {
		text = text + strconv.Itoa(i+1) + ". <b>" + link.Title + "</b>\n" + link.URL + "\n"
	}

	msgConfig := h.tg.CreateNewMessage(setting.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		return e.Wrap("error while sending response with latestLinks", err)
	}
	return nil
}

func (h LinksHandler) askLinksToken(setting *settings.Setting) error {
	err := h.settingsService.ChangeMode(setting.Username, LinksToken)
	if err != nil {
		return e.Wrap("failed to change mode in askLinksToken", err)
	}
	if err := h.tg.SendMessage(setting.ChatId, AskForLinksToken); err != nil {
		return e.Wrap("error while sending response with links", err)
	}
	return nil
}
