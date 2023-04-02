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
	event := GenerateTestMessage("/start")
	err := processor.Process(event)
	if err != nil {
		t.Fatal("There was an error during process the message" + err.Error())
	}
	settings, err := processor.SettingsService.GetByUserName("Test")
	if errors.Is(sql.ErrNoRows, err) {
		t.Error("settings did not created!")
	}
	if settings.Username != "Test" {
		t.Errorf("Expected username to be Test, got %s", settings.Username)
	}
	if settings.ChatId != 123 {
		t.Errorf("Expected username to be 123, got %d", settings.ChatId)
	}
	if processor.Tg.GetMessage() != strategies.MsgHelp {
		t.Errorf("sent messages does not match. Expected %s, got %s", strategies.MsgHelp, processor.Tg.GetMessage())
	}
}
