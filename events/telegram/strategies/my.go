package strategies

import (
	"bot/events"
	"bot/lib/e"
	"bot/settings"
	"fmt"
)

const MyCmd = "my"

const MsgMy = "Ниже информация, которая есть о вас в нашей базе 👇 \n"

type MyHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
}

func NewMyHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client) MyHandler {
	return MyHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
	}
}

func (h MyHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == MyCmd
}

func (h MyHandler) Handle(text string, setting *settings.Setting) error {
	setting, err := h.settingsService.GetByUserName(setting.Username)

	if err != nil {
		h.tg.SendMessage(setting.ChatId, NoSettingsError)
		return e.Wrap("Failed to get settings", err)
	}

	var msg = MsgMy
	msg += fmt.Sprintf("ID: %d \n", setting.ID)
	msg += fmt.Sprintf("Username: %s \n", setting.Username)
	msg += fmt.Sprintf("Chat ID: %d \n", setting.ChatId)
	msg += fmt.Sprintf("LinkAce Token: %s \n", setting.LinkaceToken)
	msg += fmt.Sprintf("EasyList Token: %s \n", setting.EasylistToken)
	msg += fmt.Sprintf("EasyWords Token: %s \n", setting.EasywordsToken)
	msg += fmt.Sprintf("Mode: %d \n", setting.Mode)
	msg += fmt.Sprintf("Context: %s \n", setting.Context)
	msg += fmt.Sprintf("EasyList ID: %d \n", setting.EasylistId)

	return h.tg.SendMessage(setting.ChatId, msg)
}
