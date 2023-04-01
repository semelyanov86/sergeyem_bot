package strategies

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/settings"
)

const HelpCmd = "help"

const MsgHelp = `Этот бот может взаимодействовать с различными сервисами, такими как LinkAce, EasyList, EasyWords. 

Он может сохранять ссылки в сервис LinkAce. Также он может вывести последние сохранённые ссылки.

Чтобы сохранить ссылку, просто отправьте её мне. Ссылка должна начинаться с https:// 🔗

Чтобы получить последние сохранённые ссылки, отправьте команду /links . В качестве аргумента вы можете передать количество ссылок, которые вы хотели бы видеть в сообщении. 🌎

Чтобы получить ссылки из определённого списка, используйте команду /listlinks. Затем выберите интересующий вас список. 🌍

Чтобы посмотреть списки, введите команду /list . 📂

Вы всегда можете вернуться в главное меню при помощи команды /cancel . 🔙
`

type HelpHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              *telegram2.Client
}

func NewHelpHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg *telegram2.Client) HelpHandler {
	return HelpHandler{
		meta:            meta,
		settingsService: settingsService,
		tg:              tg,
	}
}

func (h HelpHandler) IsSupported(mode int) bool {
	if !h.meta.Message.IsCommand() {
		return false
	}
	return h.meta.Message.Command() == HelpCmd
}

func (h HelpHandler) Handle(text string, setting *settings.Setting) error {
	return h.tg.SendMessage(setting.ChatId, MsgHelp)
}
