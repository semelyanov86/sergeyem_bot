package telegram

import (
	"bot/links"
	"bot/settings"
)

type FactoryInterface interface {
	GetLinkService(config settings.Config) links.LinkService
	SetSettings(settings *settings.Setting)
}

type FactoryResolver struct {
	Setting *settings.Setting
}

func (f *FactoryResolver) GetLinkService(config settings.Config) links.LinkService {
	return links.LinkService{
		Repository: links.LinkRepository{
			Url:   config.LinksUrl,
			Token: f.Setting.LinkaceToken,
		},
		Settings: f.Setting,
		Config:   config,
	}
}

func (f *FactoryResolver) SetSettings(settings *settings.Setting) {
	f.Setting = settings
}

type TestFactoryResolver struct {
	Setting *settings.Setting
}

func (t *TestFactoryResolver) GetLinkService(config settings.Config) links.LinkService {
	return links.LinkService{
		Repository: links.MockRepository{Token: t.Setting.LinkaceToken},
		Settings:   t.Setting,
		Config:     config,
	}
}

func (t *TestFactoryResolver) SetSettings(settings *settings.Setting) {
	t.Setting = settings
}

func NewTestFactoryResolver() *TestFactoryResolver {
	return &TestFactoryResolver{Setting: &settings.Setting{}}
}
