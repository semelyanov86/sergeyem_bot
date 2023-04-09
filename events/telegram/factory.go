package telegram

import (
	"bot/links"
	"bot/lists"
	"bot/settings"
	"bot/words"
)

type FactoryInterface interface {
	GetLinkService(config settings.Config) links.LinkService
	SetSettings(settings *settings.Setting)
	GetWordsService(config settings.Config) words.WordService
	GetListService(config settings.Config) lists.ListService
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

func (f *FactoryResolver) GetListService(config settings.Config) lists.ListService {
	return lists.ListService{
		Repository: lists.ListRepository{
			Url:   config.ListsUrl,
			Token: f.Setting.EasylistToken,
		},
		Settings: f.Setting,
		Config:   config,
	}
}

func (f *FactoryResolver) GetWordsService(config settings.Config) words.WordService {
	return words.WordService{
		Repository: words.WordRepository{
			Token: f.Setting.EasywordsToken,
			Url:   config.WordsUrl,
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

func (f *TestFactoryResolver) GetListService(config settings.Config) lists.ListService {
	return lists.ListService{
		Repository: lists.MockRepository{Token: f.Setting.EasylistToken},
		Settings:   f.Setting,
		Config:     config,
	}
}

func (f *TestFactoryResolver) GetWordsService(config settings.Config) words.WordService {
	return words.WordService{
		Repository: words.MockRepository{
			Token: f.Setting.EasywordsToken,
		},
		Settings: f.Setting,
		Config:   config,
	}
}

func (t *TestFactoryResolver) SetSettings(settings *settings.Setting) {
	t.Setting = settings
}

func NewTestFactoryResolver() *TestFactoryResolver {
	return &TestFactoryResolver{Setting: &settings.Setting{}}
}
