package main

import (
	"fmt"

	"github.com/taglme/nfc-client/pkg/client"
	"github.com/taglme/nfc-client/pkg/models"
	"github.com/taglme/nfc-client/pkg/ndefconv"
)

type handler struct {
	Exit chan bool
}

func main() {
	client := client.New("127.0.0.1:3011", "en")
	adapters, err := client.Adapters.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(adapters) == 0 {
		fmt.Println("NFC adapters not found")
		return
	}

	err = client.Ws.Connect()
	if err != nil {
		fmt.Printf("Could not connect to server: %s", err)
		return
	}

	defer client.Ws.Disconnect()

	done := make(chan bool)

	handleEvents := &handler{Exit: done}

	client.Ws.OnEvent(handleEvents.EventHandler)

	_, err = client.Jobs.Add(adapters[0].AdapterID, writeURL("http://tagl.me"))
	if err != nil {
		fmt.Println("Could not add job to server")
		return
	}

	<-done
}

func writeURL(url string) models.NewJob {
	ndefMessage := []ndefconv.NdefRecord{
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeUrl,
			Data: ndefconv.NdefRecordPayloadUrl{
				Url: url,
			},
		},
	}

	operation := models.JobStep{
		Command: models.CommandWriteNdef,
		Params: models.WriteNdefParams{
			Message: ndefMessage,
		},
	}

	operation.ToResource()

	newJob := models.NewJob{
		JobName:     "Write NDEF message",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			operation.ToResource(),
		},
	}
	return newJob
}
func (h *handler) EventHandler(e models.Event) {

	switch e.Name {
	case models.EventNameJobSubmited:
		fmt.Println("New job submited to server")
	case models.EventNameJobActivated:
		fmt.Printf("Job activated. Waiting for tag...")
	case models.EventNameTagDiscovery:
		fmt.Println("tag discovered")
	case models.EventNameRunStarted:
		fmt.Printf("Run job operations. Hold the tag...")
	case models.EventNameRunSuccess:
		fmt.Println("success")
	case models.EventNameRunError:
		fmt.Println("error")
	case models.EventNameJobFinished:
		fmt.Println("Job finished")
		h.Exit <- true
	}

}
