package controllers

import (
	"bot/events"
	"bot/settings"
	"net/http"
)
import "github.com/julienschmidt/httprouter"

func Routes(cfg settings.Config, processor events.Processor) http.Handler {
	var router = httprouter.New()
	handlers := NewHandlers(cfg, processor)
	router.HandlerFunc(http.MethodPost, "/list", handlers.ReceiveListHandler)
	router.HandlerFunc(http.MethodPost, "/meeting", handlers.ReceiveMeetingHandler)

	return handlers.EnableCORS(router)
}
