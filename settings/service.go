package settings

import (
	"bot/lib/e"
	"database/sql"
	"errors"
)

type ServiceInterface interface {
	GetOrCreateSetting(userName string, userId int64) (*Setting, error)
	UpdateSetting(setting *Setting) error
	ChangeMode(userName string, mode int) error
	GetByUserName(userName string) (*Setting, error)
	SetContext(userName string, context string) error
}

type Service struct {
	Repository RepositoryInterface
}

func (s *Service) New(db *sql.DB) Service {
	return NewSettingsService(db)
}

func NewSettingsService(db *sql.DB) Service {
	return Service{
		Repository: &Repository{DB: db},
	}
}

func (s *Service) GetByUserName(userName string) (*Setting, error) {
	setting, err := s.Repository.Get(userName)
	if errors.Is(ErrNoRecord, err) {
		return nil, nil
	}
	if err != nil {
		return nil, e.Wrap("failed to get settings", err)
	}
	return setting, nil
}

func (s *Service) GetByEasyListId(userId int64) (*Setting, error) {
	setting, err := s.Repository.GetByEasyListId(userId)
	if errors.Is(ErrNoRecord, err) {
		return nil, nil
	}
	if err != nil {
		return nil, e.Wrap("failed to get settings", err)
	}
	return setting, nil
}

func (s *Service) GetOrCreateSetting(userName string, userId int64) (*Setting, error) {
	setting, err := s.Repository.Get(userName)
	if errors.Is(ErrNoRecord, err) {
		settingModel := Setting{Username: userName, ChatId: userId}
		err = s.Repository.Insert(&settingModel)
		if err != nil {
			return nil, e.Wrap("Error while inserting setting", err)
		}
		return &settingModel, nil
	}

	return setting, nil
}

func (s *Service) UpdateSetting(setting *Setting) error {
	err := s.Repository.Update(setting)
	if err != nil {
		return e.Wrap("Failed to update setting", err)
	}
	return nil
}

func (s *Service) ChangeMode(userName string, mode int) error {
	err := s.Repository.UpdateMode(userName, mode)
	if err != nil {
		return e.Wrap("Failed to update mode", err)
	}
	return nil
}

func (s *Service) SetContext(userName string, context string) error {
	err := s.Repository.UpdateContext(userName, context)
	if err != nil {
		return e.Wrap("Failed to update mode", err)
	}
	return nil
}
