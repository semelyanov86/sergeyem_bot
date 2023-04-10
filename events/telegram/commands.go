package telegram

import (
	"bot/events"
	"bot/events/telegram/strategies"
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
	if setting == nil {
		setting = &settings.Setting{
			Username: username,
			ChatId:   meta.ChatID,
		}
	}

	p.Factory.SetSettings(setting)

	if err != nil {
		return err
	}

	var handlers = [22]CommandHandlerInterface{
		strategies.NewStartHandler(meta, p.SettingsService, p.Tg),
		strategies.NewMyHandler(meta, p.SettingsService, p.Tg),
		strategies.NewHelpHandler(meta, p.SettingsService, p.Tg),
		strategies.NewCancelHandler(meta, p.SettingsService, p.Tg),
		strategies.NewLinksHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewListsHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewLinkTokenHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewLinkSaveHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewLinkStoreHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewListLinksHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewLinksFromListHandler(meta, p.SettingsService, p.Tg, p.Factory.GetLinkService(p.config)),
		strategies.NewWordTokenHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewRandomHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewSaveWordStartHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewSaveWordTranslationsHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewSaveWordLanguageHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewStoreWordHandler(meta, p.SettingsService, p.Tg, p.Factory.GetWordsService(p.config)),
		strategies.NewBuyListsHandler(meta, p.SettingsService, p.Tg, p.Factory.GetListService(p.config)),
		strategies.NewListTokenHandler(meta, p.SettingsService, p.Tg, p.Factory.GetListService(p.config)),
		strategies.NewListIdHandler(meta, p.SettingsService, p.Tg, p.Factory.GetListService(p.config)),
		strategies.NewItemsHandler(meta, p.SettingsService, p.Tg, p.Factory.GetListService(p.config)),
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
