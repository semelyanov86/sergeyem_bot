package tests

import (
	"database/sql"
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
			"../../migrations/000001_create_settings_table.down.sql",
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
