package settings

import (
	"database/sql"
	"errors"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) Get(user string) (*Setting, error) {
	stmt := `SELECT id, username, chat_id, linkace_token, easylist_token, easywords_token, mode, context, easylist_id FROM settings
    WHERE username = ?`

	row := r.DB.QueryRow(stmt, user)

	s := &Setting{}
	err := row.Scan(&s.ID, &s.Username, &s.ChatId, &s.LinkaceToken, &s.EasylistToken, &s.EasywordsToken, &s.Mode, &s.Context, &s.EasylistId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (r *Repository) GetByEasyListId(listId int64) (*Setting, error) {
	stmt := `SELECT id, username, chat_id, linkace_token, easylist_token, easywords_token, mode, context, easylist_id FROM settings
    WHERE easylist_id = ?`

	row := r.DB.QueryRow(stmt, listId)

	s := &Setting{}
	err := row.Scan(&s.ID, &s.Username, &s.ChatId, &s.LinkaceToken, &s.EasylistToken, &s.EasywordsToken, &s.Mode, &s.Context, &s.EasylistId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (r *Repository) Insert(setting *Setting) error {
	stmt := `INSERT INTO settings (username, chat_id, linkace_token, easylist_token, easywords_token, mode, context, easylist_id)
    VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(stmt, setting.Username, setting.ChatId, setting.LinkaceToken, setting.EasylistToken, setting.EasywordsToken, setting.Mode, setting.Context, setting.EasylistId)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	setting.ID = id

	return nil
}

func (r *Repository) Update(setting *Setting) error {
	stmt := "UPDATE settings SET linkace_token = ?, easylist_token = ?, easywords_token = ?, easylist_id = ? WHERE id = ?"
	_, err := r.DB.Exec(stmt, setting.LinkaceToken, setting.EasylistToken, setting.EasywordsToken, setting.EasylistId, setting.ID)
	return err
}

func (r *Repository) UpdateMode(userName string, mode int) error {
	stmt := "UPDATE settings SET `mode` = ? WHERE username = ?"
	_, err := r.DB.Exec(stmt, mode, userName)
	return err
}

func (r *Repository) UpdateContext(userName string, context string) error {
	stmt := "UPDATE settings SET `context` = ? WHERE username = ?"
	_, err := r.DB.Exec(stmt, context, userName)
	return err
}
