package controllers

import (
	"bot/events"
	"bot/events/telegram"
	"bot/lib/e"
	"bot/settings"
	"encoding/json"
	"errors"
	"net/http"
)

const MsgListRequest = "Товары из списка "

type ListRequestData struct {
	Name   string `json:"name"`
	UserId int64  `json:"user_id"`
	ListId int    `json:"list_id"`
}

type Handlers struct {
	config settings.Config
	tg     events.Client
}

var ErrValidation = errors.New("there was a validation error, please pass correct data")

func (h Handlers) ReceiveListHandler(w http.ResponseWriter, r *http.Request) {
	var data ListRequestData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		e.ErrorResponse(w, r, http.StatusUnprocessableEntity, ErrValidation)
		return
	}

	if data.Name == "" || data.ListId < 1 || data.UserId < 1 {
		e.ErrorResponse(w, r, http.StatusUnprocessableEntity, ErrValidation)
		return
	}

	settingsService := settings.NewSettingsService(h.config.Db.Sql)
	settingsModel, err := settingsService.GetByEasyListId(data.UserId)
	if err != nil {
		e.ErrorResponse(w, r, 500, err)
		return
	}
	factory := telegram.FactoryResolver{Setting: settingsModel}
	listService := factory.GetListService(h.config)
	items, err := listService.GetItemsFromList(data.ListId)
	if err != nil {
		e.ErrorResponse(w, r, 500, err)
		return
	}
	text := MsgListRequest + data.Name + "\n" + listService.GenerateMessageFromItems(items)
	msgConfig := h.tg.CreateNewMessage(settingsModel.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		e.ErrorResponse(w, r, 500, err)
		return
	}
}

func NewHandlers(config settings.Config, processor events.Processor) Handlers {
	return Handlers{
		config: config,
		tg:     processor.GetTgClient(),
	}
}
