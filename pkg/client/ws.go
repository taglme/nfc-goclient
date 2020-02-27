package client

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/taglme/nfc-client/pkg/models"
)

type WsService interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	OnEvent(func(handler EventHandler))
	SetLocale(locale string) error
	ConnString() string
}

type EventHandler func(e models.Event)
type wsService struct {
	conn     *websocket.Conn
	handlers []EventHandler
}

func (ws *wsService) OnEvent(handler EventHandler) {
	ws.handlers = append(ws.handlers, handler)
}

func (ws *wsService) eventListener(e models.Event) {
	for _, handler := range ws.handlers {
		handler(e)
	}
}

func (ws *wsService) read() {
	defer ws.conn.Close()
	for {
		var eventResource models.EventResource
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			break
		}
		err = json.Unmarshal(message, &eventResource)
		if err != nil {
			//log here
			continue
		}

		event, err := eventResource.ToEvent()
		if err != nil {
			//log here
			continue
		}
		ws.eventListener(event)

	}

}
