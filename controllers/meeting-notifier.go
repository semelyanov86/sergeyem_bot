package controllers

import (
	"bot/lib/e"
	"bot/settings"
	"encoding/json"
	"net/http"
)

const AdminUserName = "sergeyem"

const MsgMeetingRequest = "Назначена новая встреча в сервисе EasyAppointments:\n"

type MeetingData struct {
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
	Email         string `json:"email"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Timezone      string `json:"timezone"`
}

func (h Handlers) ReceiveMeetingHandler(w http.ResponseWriter, r *http.Request) {
	var data MeetingData

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		e.ErrorResponse(w, r, http.StatusUnprocessableEntity, ErrValidation)
		return
	}

	if data.StartDatetime == "" || data.EndDatetime == "" || data.Email == "" {
		e.ErrorResponse(w, r, http.StatusUnprocessableEntity, ErrValidation)
		return
	}

	settingsService := settings.NewSettingsService(h.config.Db.Sql)
	settingsModel, err := settingsService.GetByUserName(AdminUserName)
	if err != nil {
		e.ErrorResponse(w, r, 500, err)
		return
	}

	text := MsgMeetingRequest + "Дата начала: " + data.StartDatetime + "\nДата окончания: " + data.EndDatetime + "\nИмя: " + data.FirstName + "\nФамилия: " + data.LastName + "\nВременная зона: " + data.Timezone + "\nПроверьте свой календарь и добавьте ссылку на Google Meet. Клиент ждёт! 🍾"
	msgConfig := h.tg.CreateNewMessage(settingsModel.ChatId, text)
	if err := h.tg.Send(msgConfig); err != nil {
		e.ErrorResponse(w, r, 500, err)
		return
	}
}
