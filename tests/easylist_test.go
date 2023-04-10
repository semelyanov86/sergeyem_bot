package tests

import (
	"bot/events/telegram/strategies"
	"bot/lists"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"testing"
)

func TestGettingListsAskForToken(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithoutTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.BuyListsCmd)
	err = processor.Process(event)
	if err != nil {
		t.Error("There was an error during request - " + err.Error())
	}
	if processor.Tg.GetMessage() != strategies.AskForListsToken {
		t.Errorf("result message does not match. Expected %s, got %s", strategies.AskForListsToken, processor.Tg.GetMessage())
	}
	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.Mode != strategies.AskedEasylistToken {
		t.Errorf("expected mode to be %d, got %d", strategies.AskedEasylistToken, setting.Mode)
	}
}

func TestSaveEasyListToken(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithoutTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskedEasylistToken)
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("NEW_TOKEN")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process original entry" + err.Error())
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.EasylistToken != "NEW_TOKEN" {
		t.Errorf("expected token to be %s, got %s", "NEW_TOKEN", setting.EasylistToken)
	}
	if processor.Tg.GetMessage() != strategies.MsgAskEasyListId {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgAskEasyListId)
	}
}

func TestSaveEasylistId(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskedEasylistId)
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("11")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process original entry" + err.Error())
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.EasylistId != 11 {
		t.Errorf("expected token to be %s, got %s", "NEW_TOKEN", setting.EasylistToken)
	}
	if processor.Tg.GetMessage() != strategies.MsgListId {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgListId)
	}
}

func TestGettingListsFromEasyList(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.BuyListsCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process lists from Easylist" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgBuyLists) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgBuyLists)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "Soupe Schi") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Soupe Schi")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "Borsch") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Borsch")
	}
}

func TestGettingItemsWithoutListThroghsError(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.ItemsCmd)
	err = processor.Process(event)
	if err == nil {
		t.Error("There should be wrong list error, got null")
	}

	if processor.Tg.GetMessage() != strategies.MsgErrorWrongList {
		t.Errorf("sent messages does not match. Expected that message %s equals %s", processor.Tg.GetMessage(), strategies.MsgErrorWrongList)
	}
}

func TestGettingItemsFromSpecificList(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.ItemsCmd + " 3")
	event.Meta.Message.Entities[0].Length = 6
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process items from Easylist" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), lists.MsgItems) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), lists.MsgItems)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "<b>Kapusta</b> (1 st)") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Kapusta (1 st)")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "<b>Svekla</b> (2 st)") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Svekla (2 st)")
	}
}
