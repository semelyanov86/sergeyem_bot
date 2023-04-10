package settings

import (
	"database/sql"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Port          int
	Env           string
	Db            *Database
	AppName       string `yaml:"appName"`
	LinksUrl      string `yaml:"linksUrl"`
	LinksPerPage  int    `yaml:"linksPerPage"`
	WordsUrl      string `yaml:"wordsUrl"`
	WordsPerPage  int    `yaml:"wordsPerPage"`
	WordsLanguage string `yaml:"wordsLanguage"`
	ListsUrl      string `yaml:"listsUrl"`
	ListsPerPage  int    `yaml:"listsPerPage"`
	WebhookUrl    string `yaml:"webhookUrl"`
}

type Database struct {
	Dsn          string
	Host         string
	Login        string
	Password     string
	Dbname       string
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxIdleTime  string `yaml:"maxIdleTime"`
	Sql          *sql.DB
}

func ReadConfigFile(cfg *Config) {
	dirname, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	bytesOut, err := os.ReadFile(dirname + "/chatbot.yaml")

	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(bytesOut, &cfg); err != nil {
		panic(err)
	}
}
