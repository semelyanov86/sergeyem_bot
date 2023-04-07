package tests

import (
	"bot/events/telegram/strategies"
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
	if err == nil {
		t.Error("There should be a token error, got nil")
	}
	if processor.Tg.GetMessage() != strategies.AskForWordsToken {
		t.Errorf("result message does not match. Expected %s, got %s", strategies.AskForWordsToken, processor.Tg.GetMessage())
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
