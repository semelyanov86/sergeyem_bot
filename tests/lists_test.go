package tests

import (
	"bot/events/telegram/strategies"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"testing"
)

func TestGettingAllLists(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.ListsCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process lists message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgLists) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgLinks)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "PHP") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "PHP")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "Golang") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Golang")
	}
}

func TestLinksFromListCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.ListLinksCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process listslinks message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgListLinks) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgListLinks)
	}
}

func TestDisplayingLinksFromList(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskListForLinks)
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("there was an error during receiving setting" + err.Error())
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("PHP|2")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process listslinks message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgLinksFromList) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgLinksFromList)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "https://sergeyem.ru") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "https://sergeyem.ru")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "https://sergeyem.ru") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "https://itvolga.com")
	}
}
