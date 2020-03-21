package client

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/taglme/nfc-goclient/pkg/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		//_, _, err := c.ReadMessage()
		//if err != nil {
		//	log.Printf(err.Error())
		//	break
		//}
		resp, err := json.Marshal(models.EventResource{
			EventID:     "123",
			Name:        models.EventNameAdapterDiscovery.String(),
			AdapterID:   "123",
			AdapterName: "aname",
			Data:        nil,
			CreatedAt:   "2006-01-02T15:04:05Z",
		})

		if err != nil {
			log.Fatal("Can't marshall test model\n", err)
		}

		err = c.WriteMessage(websocket.TextMessage, resp)
		if err != nil {
			break
		}
	}
}

func echoForErr(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		//_, _, err := c.ReadMessage()
		//if err != nil {
		//	log.Printf(err.Error())
		//	break
		//}
		err = c.WriteMessage(websocket.TextMessage, []byte("Different model"))
		if err != nil {
			break
		}
	}
}

func echoLocales(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf(string(msg), err.Error())
			break
		}
		//err = c.WriteMessage(websocket.TextMessage, []byte("Different model"))
		//if err != nil {
		//	break
		//}
	}
}

func TestWsService_ConnString(t *testing.T) {
	s := newWsService("url")
	url := s.ConnString()
	assert.Equal(t, "url/ws", url)
}

func TestWsService_Connect(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws := newWsService(u)

	assert.Equal(t, false, ws.IsConnected())
	err := ws.Connect()
	assert.Nil(t, err)
	assert.Equal(t, true, ws.IsConnected())
	err = ws.Disconnect()
	assert.Nil(t, err)
	assert.Equal(t, false, ws.IsConnected())
	//ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	//if err != nil {
	//	t.Fatalf("%v", err)
	//}
	//defer ws.Close()
	//
	//// Send message to server, read response and check to see if it's what we expect.
	//for i := 0; i < 10; i++ {
	//	if err := ws.WriteMessage(websocket.TextMessage, []byte("hello")); err != nil {
	//		t.Fatalf("%v", err)
	//	}
	//	_, p, err := ws.ReadMessage()
	//	if err != nil {
	//		t.Fatalf("%v", err)
	//	}
	//	if string(p) != "hello" {
	//		t.Fatalf("bad message")
	//	}
	//}
}

func TestWsService_OnEvent(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws := newWsService(u)
	err := ws.Connect()
	if err != nil {
		log.Fatal("Can't connect to test server\n")
	}

	a := 0
	l := func(e models.Event) {
		// change value to validate if handler is working
		assert.Equal(t, models.EventNameAdapterDiscovery, e.Name)
		a++
	}

	ws.OnEvent(l)

	// sleep to let handler work
	time.Sleep(time.Second)

	assert.Less(t, 0, a)
}

func TestWsService_OnError(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echoForErr))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws := newWsService(u)
	err := ws.Connect()
	if err != nil {
		log.Fatal("Can't connect to test server\n")
	}

	a := 0
	l := func(e error) {
		// change value to validate if handler is working
		//assert.Equal(t, models.EventNameAdapterDiscovery, e.Name)
		assert.EqualError(t, e, "Can't unmarshall event resource\n: invalid character 'D' looking for beginning of value")
		a++
	}

	ws.OnError(l)

	// sleep to let handler work
	time.Sleep(time.Second)

	assert.Less(t, 0, a)
}

func TestWsService_SetLocale(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echoLocales))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws := newWsService(u)

	err := ws.SetLocale("en")
	assert.EqualError(t, err, "Can't set locale. Connection were not initialized")

	err = ws.Connect()
	if err != nil {
		log.Fatal("Can't connect to test server\n")
	}

	err = ws.SetLocale("en")

	assert.Nil(t, err)
}
