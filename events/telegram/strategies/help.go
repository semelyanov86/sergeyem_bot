package strategies

import (
	"bot/events"
	"bot/settings"
)

const HelpCmd = "help"

const MsgHelp = `Этот бот может взаимодействовать с различными сервисами, такими как LinkAce, EasyList, EasyWords. 

Он может сохранять ссылки в сервис LinkAce. Также он может вывести последние сохранённые ссылки.

Чтобы сохранить ссылку, просто отправьте её мне. Ссылка должна начинаться с https:// 🔗

Чтобы получить последние сохранённые ссылки, отправьте команду /links . В качестве аргумента вы можете передать количество ссылок, которые вы хотели бы видеть в сообщении. 🌎

Чтобы получить ссылки из определённого списка, используйте команду /listlinks. Затем выберите интересующий вас список. 🌍

Чтобы посмотреть списки, введите команду /lists . 📂

Хотите начать учить новые слова? Отправьте команду /random, которая также принимает в качестве параметра количество слов для изучения. Система получит слова в рандомном порядке из сервиса EasyWords. Если в настройках EasyWords указано, что изученные слова скрываются, их в выдаче не будет.

Чтобы сохранить слово, наберите команду /saveword . После этого мы спросим вас оригинальное значение слова, его перевод и язык.

Для получения информации по вашему пользователю (часто нужно для интеграции со сторонними системами), отправьте команду /my 😎

Для получения всех список сервиса EasyList, нужно передать команду /buylist , которая принимает в качестве параметра количество элементов в сообщении

Для получения списка покупок из списка, вызовите команду /items , в которую необходимо передать ID списка

Вы всегда можете вернуться в главное меню при помощи команды /cancel . 🔙
`

type HelpHandler struct {
	meta            events.TelegramMeta
	settingsService settings.ServiceInterface
	tg              events.Client
}

func NewHelpHandler(meta events.TelegramMeta, settingsService settings.ServiceInterface, tg events.Client) HelpHandler {
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
