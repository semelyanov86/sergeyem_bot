package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/lib/e"
	"bot/links"
	"bot/settings"
)

const MsgStoreLink = "Ссылка успешно сохранена! 💋"

const WaitingMsg = "Начинаем сохранять вашу ссылку. На это может потребоваться больше 20 секунд. Ожидайте... 👷‍♀️"

const LinkSaveError = "Нельзя сохранить ссылку. Произошла странная ошибка - "

type LinkStoreHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
	linkService     links.LinkService
}

func NewLinkStoreHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client, linkService links.LinkService) LinkStoreHandler {
	return LinkStoreHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
		linkService:     linkService,
	}
}

func (h LinkStoreHandler) IsSupported(mode int) bool {
	if h.meta.Message.IsCommand() {
		return false
	}
	return mode == AskList
}

func (h LinkStoreHandler) Handle(msg string, setting *settings.Setting) error {
	listsModel := make([]links.List, 1)
	listsModel[0] = links.List{Name: msg}

	linkModel := links.Link{
		URL:   setting.Context,
		Lists: listsModel,
	}
	err := h.settingsService.ChangeMode(setting.Username, Root)
	if err != nil {
		_ = h.tg.SendMessage(setting.ChatId, ModeChangeFailed)
		return e.Wrap("failed to change mode to AskList", err)
	}
	err = h.settingsService.SetContext(setting.Username, "")
	if err != nil {
		h.tg.SendMessage(setting.ChatId, ContextChangeFailed)
		return e.Wrap("failed to change context", err)
	}

	_ = h.tg.SendMessage(setting.ChatId, WaitingMsg)
	err = h.linkService.SaveLink(&linkModel)
	if err != nil {
		h.tg.SendMessage(setting.ChatId, LinkSaveError+err.Error())
		return e.Wrap("failed to save link", err)
	}

	return h.tg.SendMessage(setting.ChatId, MsgStoreLink)
}
