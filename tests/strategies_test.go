package tests

import (
	"bot/events/telegram/strategies"
	"database/sql"
	"errors"
	"testing"
)

func TestStartStrategy(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	event := GenerateTestMessage("/" + strategies.StartCmd)
	err := processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process the message" + err.Error())
	}
	settings, err := processor.SettingsService.GetByUserName(TestUserName)
	if errors.Is(sql.ErrNoRows, err) {
		t.Error("settings did not created!")
	}
	if settings.Username != TestUserName {
		t.Errorf("Expected username to be Test, got %s", settings.Username)
	}
	if settings.ChatId != TestChatId {
		t.Errorf("Expected username to be 123, got %d", settings.ChatId)
	}
	if processor.Tg.GetMessage() != strategies.MsgHelp {
		t.Errorf("sent messages does not match. Expected %s, got %s", strategies.MsgHelp, processor.Tg.GetMessage())
	}
}

func TestHelpStrategy(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	event := GenerateTestMessage("/" + strategies.HelpCmd)
	err := processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process help message" + err.Error())
	}
	if processor.Tg.GetMessage() != strategies.MsgHelp {
		t.Errorf("sent messages does not match. Expected %s, got %s", strategies.MsgHelp, processor.Tg.GetMessage())
	}
}

func TestCancelStrategy(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	_, err := processor.SettingsService.GetOrCreateSetting(TestUserName, TestChatId)
	if err != nil {
		t.Fatal("there was an error during creation of setting" + err.Error())
	}
	processor.SettingsService.SetContext(TestUserName, "some_context")
	err = processor.SettingsService.ChangeMode(TestUserName, strategies.AskList)
	if err != nil {
		t.Fatal("there was an error during change of mode" + err.Error())
	}
	event := GenerateTestMessage("/" + strategies.CancelCmd)
	err = processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process cancel message" + err.Error())
	}
	if processor.Tg.GetMessage() != strategies.MsgCancel {
		t.Errorf("sent messages does not match. Expected %s, got %s", strategies.MsgHelp, processor.Tg.GetMessage())
	}
	setting, err := processor.SettingsService.GetByUserName(TestUserName)
	if err != nil {
		t.Fatal("There was an error while getting setting by user" + err.Error())
	}
	if setting.Context != "" {
		t.Errorf("expected empty context, got %s", setting.Context)
	}
	if setting.Mode != strategies.Root {
		t.Errorf("expected mode to be %d, got %d", strategies.Root, setting.Mode)
	}
}

func TestDefaultStrategy(t *testing.T) {
	processor, down := NewTestProcessorWithDb(t)
	defer down()
	event := GenerateTestMessage("/some_unknown_command")
	err := processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process help message" + err.Error())
	}
	if processor.Tg.GetMessage() != strategies.MsgDefault {
		t.Errorf("sent messages does not match. Expected %s, got %s", strategies.MsgDefault, processor.Tg.GetMessage())
	}
}
