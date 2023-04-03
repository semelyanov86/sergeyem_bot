package tests

import (
	"bot/events/telegram/strategies"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/octoper/go-ray"
	"strings"
	"testing"
)

func TestGettingListAskForToken(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithoutTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.LinksCmd)
	err = processor.Process(event)
	if err == nil {
		t.Error("There should be a token error, got nil")
	}
	if processor.Tg.GetMessage() != strategies.AskForLinksToken {
		t.Errorf("result message does not match. Expected %s, got %s", strategies.AskForLinksToken, processor.Tg.GetMessage())
	}
}

func TestGetAllLinks(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.LinksCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgLinks) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgLinks)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "https://sergeyem.ru") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "https://sergeyem.ru")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "https://itvolga.com") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "https://itvolga.com")
	}
}

func TestPaginationFilter(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.LinksCmd + " 1")
	event.Meta.Message.Entities[0].Length = 6
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgLinks) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgLinks)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "https://sergeyem.ru") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "https://sergeyem.ru")
	}
	if strings.Contains(processor.Tg.GetMessage(), "https://itvolga.com") {
		t.Errorf("Message contains link https://itvolga.com but this link should not exist")
	}
}

func TestLinkTokenSavedSuccessfully(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithoutTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.LinksToken)
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("SOME_TOKEN_NEW")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgLinkToken {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgLinkToken)
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.LinkaceToken != "SOME_TOKEN_NEW" {
		t.Errorf("Token does not saved, expected %s got %s", "SOME_TOKEN_NEW", setting.LinkaceToken)
	}
}

func TestSendingLinkAskForLists(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("https://sergeyem.ru")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgSaveLink) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgSaveLink)
	}
}

func TestSendingWrongLinkIsNotFoundCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("sergeyem.ru")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgDefault {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgDefault)
	}
}

func TestStoringLink(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskList)
	processor.SettingsService.SetContext(TestUserName, "https://sergeyem.ru")
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("there was an error during getting new setting" + err.Error())
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("PHP")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process links message" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgStoreLink {
		t.Errorf("sent messages does not match. Expected that message %s equals %s", processor.Tg.GetMessage(), strategies.MsgStoreLink)
	}
}
