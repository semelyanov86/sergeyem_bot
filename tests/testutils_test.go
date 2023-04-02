package tests

import (
	telegram2 "bot/clients/telegram"
	"bot/events"
	"bot/events/telegram"
	"bot/settings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"testing"
)

func NewTestDB(t *testing.T) (*sql.DB, func()) {
	var dbCred = os.Getenv("BOT_TEST_DB")
	db, err := sql.Open("mysql", dbCred)
	if err != nil {
		t.Fatal(err)
	}

	migrations := [...]string{
		"../migrations/000001_create_settings_table.up.sql",
	}
	for _, migration := range migrations {
		script, err := os.ReadFile(migration)
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	}
	return db, func() {
		migrations := [...]string{
			"../migrations/000001_create_settings_table.down.sql",
		}
		for _, migration := range migrations {
			script, err := os.ReadFile(migration)
			if err != nil {
				t.Fatal(err)
			}
			_, err = db.Exec(string(script))
			if err != nil {
				t.Fatal(err)
			}
		}
		db.Close()
	}
}

func NewTestProcessorWithDb(t *testing.T) (*telegram.Processor, func()) {
	db, teardown := NewTestDB(t)
	var testClient = telegram2.NewTestClient()
	var cfg settings.Config
	cfg.Db = &settings.Database{}
	cfg.Db.Dsn = os.Getenv("BOT_TEST_DB")
	cfg.Db.Sql = db
	return telegram.New(testClient, cfg), teardown
}

func GenerateTestMessage(text string) events.Event[events.TelegramMeta] {
	return events.Event[events.TelegramMeta]{
		Type: events.Message,
		Text: text,
		Meta: events.TelegramMeta{
			ChatID:   123,
			Username: "Test",
			Message: &tgbotapi.Message{
				Text: text,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Length: len(text),
					},
				},
			},
		},
	}
}
