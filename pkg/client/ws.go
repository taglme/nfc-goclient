package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/taglme/nfc-client/pkg/models"
)

type EventHandler func(e models.Event)
type ErrorHandler func(e error)

type WsService interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	OnEvent(EventHandler)
	ConnString() string
	OnError(ErrorHandler)
	SetLocale(locale string) error
}

type wsService struct {
	url           string
	path          string
	conn          *websocket.Conn
	handlers      []EventHandler
	errorHandlers []ErrorHandler
}

func newWsService(url string) WsService {
	return &wsService{
		url:           url,
		path:          "/ws",
		handlers:      make([]EventHandler, 0),
		errorHandlers: make([]ErrorHandler, 0),
	}
}

func (s *wsService) Connect() (err error) {
	s.conn, _, err = websocket.DefaultDialer.Dial(s.url+s.path, nil)
	if err != nil {
		return errors.Wrap(err, "Can't connect to the ws endpoint\n")
	}

	go s.read()
	return nil
}

func (s *wsService) IsConnected() bool {
	return s.conn != nil
}

func (s *wsService) Disconnect() (err error) {
	err = s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "WS connection closed"))
	if err != nil {
		return errors.Wrap(err, "Error on close WS connection")
	}
	s.conn = nil
	return nil
}

func (s *wsService) SetLocale(locale string) (err error) {
	if s.conn == nil {
		return errors.New("Can't set locale. Connection were not initialized")
	}

	body, err := json.Marshal(models.SetLocaleParamsResource{Locale: locale})
	if err != nil {
		return errors.Wrap(err, "Error on marshall set locale resource")
	}

	err = s.conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return errors.Wrap(err, "Error on send set locale resource")
	}

	return nil
}

func (s *wsService) OnEvent(handler EventHandler) {
	s.handlers = append(s.handlers, handler)
}

func (s *wsService) OnError(h ErrorHandler) {
	s.errorHandlers = append(s.errorHandlers, h)
}

func (s *wsService) eventListener(e models.Event) {
	for _, handler := range s.handlers {
		handler(e)
	}
}

func (s *wsService) errListener(e error) {
	for _, handler := range s.errorHandlers {
		handler(e)
	}
}

func (s *wsService) read() {
	defer func() {
		s.conn.Close()
		s.conn = nil
	}()

	for {
		var eventResource models.EventResource
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't read WS message\n"))
			return
		}
		err = json.Unmarshal(message, &eventResource)
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't unmarshall event resource\n"))
			continue
		}

		event, err := eventResource.ToEvent()
		if err != nil {
			s.errListener(errors.Wrap(err, "Can't convert event resource to the event model\n"))

			continue
		}
		s.eventListener(event)
	}
}

func (s *wsService) ConnString() string {
	return s.url + s.path
}
