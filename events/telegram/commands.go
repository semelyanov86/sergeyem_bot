package telegram

import (
	"bot/events"
	"bot/events/telegram/strategies"
	"bot/links"
	"bot/settings"
	"strings"
)

type CommandHandlerInterface interface {
	IsSupported(mode int) bool
	Handle(text string, setting *settings.Setting) error
}

func (p *Processor) doCmd(text string, meta events.TelegramMeta) error {
	username := meta.Username
	setting, err := p.getSettingsModel(username)
	if err != nil {
		return err
	}

	var handlers = [11]CommandHandlerInterface{
		strategies.NewStartHandler(meta, p.SettingsService, p.Tg),
		strategies.NewHelpHandler(meta, p.SettingsService, p.Tg),
		strategies.NewCancelHandler(meta, p.SettingsService, p.Tg),
		strategies.NewLinksHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewListsHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewLinkTokenHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewLinkSaveHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewLinkStoreHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewListLinksHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewLinksFromListHandler(meta, p.SettingsService, p.Tg, p.getLinkService(setting)),
		strategies.NewDefaultHandler(meta, p.SettingsService, p.Tg),
	}
	for _, handler := range handlers {
		if handler.IsSupported(setting.Mode) {
			err := handler.Handle(strings.TrimSpace(text), setting)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (p *Processor) getLinkService(setting *settings.Setting) links.LinkService {
	return links.LinkService{
		Repository: links.LinkRepository{
			Url:   p.config.LinksUrl,
			Token: setting.LinkaceToken,
		},
		Settings: setting,
		Config:   p.config,
	}
}
