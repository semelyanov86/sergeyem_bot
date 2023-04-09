package tests

import (
	"bot/events/telegram/strategies"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
	"testing"
)

func TestGettingRandomWordAskForToken(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithoutTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.RandomCmd)
	err = processor.Process(event)
	if err != nil {
		t.Error(err.Error())
	}
	if processor.Tg.GetMessage() != strategies.AskForWordsToken {
		t.Errorf("result message does not match. Expected %s, got %s", strategies.AskForWordsToken, processor.Tg.GetMessage())
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.Mode != strategies.AskWordsToken {
		t.Errorf("expected mode to be %d, got %d", strategies.AskWordsToken, setting.Mode)
	}
}

func TestGetRandomWords(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.RandomCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process words message" + err.Error())
	}

	if !strings.Contains(processor.Tg.GetMessage(), strategies.MsgRandom) {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), strategies.MsgRandom)
	}
	if !strings.Contains(processor.Tg.GetMessage(), "Hamburg") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Hamburg")
	}
	if !strings.Contains(processor.Tg.GetMessage(), "Stabil") {
		t.Errorf("sent messages does not match. Expected that message %s contains %s", processor.Tg.GetMessage(), "Stabil")
	}
}

func TestStartSavingWordCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	setting, err := GenerateTestUserWithTokens(processor)
	processor.Factory.SetSettings(setting)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}

	event := GenerateTestMessage("/" + strategies.SaveWordStartCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during start save word command" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgSaveStart {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgSaveStart)
	}

	setting, _ = processor.SettingsService.GetByUserName(TestUserName)
	if setting.Mode != strategies.AskedWordOriginal {
		t.Errorf("mode does not match, expected %d, got %d", setting.Mode, strategies.AskedWordOriginal)
	}
}

func TestAskForOriginalCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskedWordOriginal)
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("Moin")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process original entry" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgAskTranslation {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgAskTranslation)
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.Mode != strategies.AskedWordTranslation {
		t.Errorf("mode does not match, expected %d, got %d", setting.Mode, strategies.AskedWordOriginal)
	}
	if setting.Context != "Moin" {
		t.Errorf("expected context to be %s, got %s", "Moin", setting.Context)
	}
}

func TestAskForTranslationCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskedWordTranslation)
	processor.SettingsService.SetContext(TestUserName, "Moin")
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("Translated")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process original entry" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgAskWordLanguage {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgAskWordLanguage)
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.Mode != strategies.AskedWordLanguage {
		t.Errorf("mode does not match, expected %d, got %d", setting.Mode, strategies.AskedWordLanguage)
	}
	if setting.Context != "Moin|Translated" {
		t.Errorf("expected context to be %s, got %s", "Moin|Translated", setting.Context)
	}
}

func TestStoreWordCommand(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := GenerateTestUserWithTokens(processor)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.ChangeMode(TestUserName, strategies.AskedWordLanguage)
	processor.SettingsService.SetContext(TestUserName, "Moin|Translated")
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	processor.Factory.SetSettings(setting)

	event := GenerateTestMessage("DE")
	event.Meta.Message.Entities[0] = tgbotapi.MessageEntity{}
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process original entry" + err.Error())
	}

	if processor.Tg.GetMessage() != strategies.MsgStoreWord {
		t.Errorf("sent messages does not match. Expected that message %s equals to %s", processor.Tg.GetMessage(), strategies.MsgStoreWord)
	}

	setting, err = processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("Failed to get updated settings")
	}
	if setting.Mode != strategies.Root {
		t.Errorf("mode does not match, expected %d, got %d", setting.Mode, strategies.Root)
	}
	if setting.Context != "" {
		t.Errorf("expected context to be %s, got %s", "", setting.Context)
	}
}
