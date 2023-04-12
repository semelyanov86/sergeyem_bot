package main

import (
	"bot/clients/telegram"
	event_consumer "bot/consumer/event-consumer"
	telegram2 "bot/events/telegram"
	"bot/settings"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

const tgBotHost = "api.telegram.org"

const BatchSize = 100

func main() {
	var cfg settings.Config

	settings.ReadConfigFile(&cfg)
	cfg.Db.Dsn = cfg.Db.Login + ":" + cfg.Db.Password + "@" + cfg.Db.Host + "/" + cfg.Db.Dbname + "?parseTime=true"

	db, err := openDB(cfg.Db.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	cfg.Db.Sql = db

	tgClient := telegram.New(tgBotHost, mustToken())
	eventsProcessor := telegram2.New(tgClient, cfg)

	log.Println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, BatchSize, cfg)

	if err := consumer.Start(cfg.WebhookUrl, cfg.Port); err != nil {
		log.Fatal("service is stopped", err)
	}
	log.Println("we are done")
}

func mustToken() string {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Telegram token is not specified")
	}
	return token
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
