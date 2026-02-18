package main

import (
	"fmt"

	"github.com/taglme/nfc-goclient/pkg/client"
	"github.com/taglme/nfc-goclient/pkg/models"
	"github.com/taglme/nfc-goclient/pkg/ndefconv"
)

type Tasker struct {
	client *client.Client
	tasks  []models.NewJob
	done   chan bool
}

func NewTasker(baseURL string, xAppKey string, xUserToken string, locale string) *Tasker {
	return &Tasker{
		client: client.NewWithOptions(client.Options{
			BaseURL:    baseURL,
			Locale:     locale,
			XAppKey:    xAppKey,
			XUserToken: xUserToken,
		}),
		tasks: make([]models.NewJob, 0),
		done:  make(chan bool),
	}
}
func (t *Tasker) Add(nj models.NewJob) {
	t.tasks = append(t.tasks, nj)
}

func (t *Tasker) Run(adapterNumber int) error {
	err := t.client.Ws.Connect()
	if err != nil {
		return err
	}

	info, err := t.client.About.Get()
	if err != nil {
		return err
	}

	printInfo(info)
	fmt.Printf("Success connection to host\n")
	defer t.client.Ws.Disconnect()

	t.client.Ws.OnEvent(t.EventHandler)

	adapters, err := t.client.Adapters.GetAll()
	if err != nil {
		return err
	}

	printAdapters(adapters)

	if len(adapters) == 0 {
		return fmt.Errorf("Adapters not found")
	}

	for _, task := range t.tasks {
		fmt.Printf(" ---TEST %s BEGIN---\n", task.JobName)
		_, err := t.client.Jobs.Add(adapters[adapterNumber-1].AdapterID, task)
		if err != nil {
			return err
		}

		<-t.done
		fmt.Printf(" ---TEST %s FINISHED---\n", task.JobName)

	}
	return nil
}
func (t *Tasker) EventHandler(e models.Event) {
	switch e.Name {
	case models.EventNameJobSubmited:
		fmt.Println("[Job submitted event]")
		job, ok := e.GetJob()
		if !ok {
			fmt.Println("Job not found in event data")
		} else {
			fmt.Println(job)
		}

	case models.EventNameJobActivated:
		fmt.Println("[Job activated event]")
		job, ok := e.GetJob()
		if !ok {
			fmt.Println("Job not found in event data")
		} else {
			fmt.Println(job)
		}

	case models.EventNameTagDiscovery:
		fmt.Println("[Tag discovered event]")
		tag, ok := e.GetTag()
		if !ok {
			fmt.Println("Tag not found in event data")
		} else {
			fmt.Println(tag)
		}

	case models.EventNameTagRelease:
		fmt.Println("[Tag release event]")
		tag, ok := e.GetTag()
		if !ok {
			fmt.Println("Tag not found in event data")
		} else {
			fmt.Println(tag)
		}

	case models.EventNameAdapterDiscovery:
		fmt.Println("[Adapter discovery event]")
		adapter, ok := e.GetAdapter()
		if !ok {
			fmt.Println("Adapter not found in event data")
		} else {
			fmt.Println(adapter)
		}

	case models.EventNameAdapterRelease:
		fmt.Println("[Adapter release event]")
		adapter, ok := e.GetAdapter()
		if !ok {
			fmt.Println("Adapter not found in event data")
		} else {
			fmt.Println(adapter)
		}

	case models.EventNameRunStarted:
		fmt.Println("[Run started event]")
		run, ok := e.GetRun()
		if !ok {
			fmt.Println("Run not found in event data")
		} else {
			fmt.Println(run)
		}
		fmt.Println(e.Data)

	case models.EventNameRunSuccess:
		fmt.Println("[Run success event]")
		run, ok := e.GetRun()
		if !ok {
			fmt.Println("Run not found in event data")
		} else {
			fmt.Println(run)
		}
		fmt.Println(e.Data)

	case models.EventNameRunError:
		fmt.Println("[Run error event]")
		run, ok := e.GetRun()
		if !ok {
			fmt.Println("Run not found in event data")
		} else {
			fmt.Println(run)
		}
		fmt.Println(e.Data)

	case models.EventNameJobFinished:
		fmt.Println("[Job finished event]")
		job, ok := e.GetJob()
		if !ok {
			fmt.Println("Job not found in event data")
		} else {
			fmt.Println(job)
		}
		t.done <- true
	}
}

func TaskWriteURL(url string) models.NewJob {
	ndefMessage := []ndefconv.NdefRecord{
		ndefconv.NdefRecord{
			Type: ndefconv.NdefRecordPayloadTypeUrl,
			Data: ndefconv.NdefRecordPayloadUrl{
				Url: url,
			},
		},
	}
	opWrite := models.JobStep{
		Command: models.CommandWriteNdef,
		Params: models.WriteNdefParams{
			Message: ndefMessage,
		},
	}
	newJob := models.NewJob{
		JobName:     "Write NDEF message",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opWrite.ToResource(),
		},
	}
	return newJob
}

func TaskRead() models.NewJob {
	opGetTags := models.JobStep{
		Command: models.CommandGetTags,
	}

	opRead := models.JobStep{
		Command: models.CommandReadNdef,
	}

	newJob := models.NewJob{
		JobName:     "Read tag",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opGetTags.ToResource(),
			opRead.ToResource(),
		},
	}

	return newJob
}

func TaskTransmit(txBytes []byte) models.NewJob {
	opTransmit := models.JobStep{
		Command: models.CommandTransmitTag,
		Params: models.TransmitTagParams{
			TxBytes: txBytes,
		},
	}
	newJob := models.NewJob{
		JobName:     "Transmit tag",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opTransmit.ToResource(),
		},
	}
	return newJob
}

func TaskGetDump() models.NewJob {
	opDump := models.JobStep{
		Command: models.CommandGetDump,
	}

	newJob := models.NewJob{
		JobName:     "Dump tag",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opDump.ToResource(),
		},
	}
	return newJob
}

func TaskFormatDefault() models.NewJob {
	opFormat := models.JobStep{
		Command: models.CommandFormatDefault,
	}

	newJob := models.NewJob{
		JobName:     "Format tag",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opFormat.ToResource(),
		},
	}
	return newJob
}

func TaskSetPassword(password []byte) models.NewJob {
	opSetPassword := models.JobStep{
		Command: models.CommandSetPassword,
		Params: models.SetPasswordParams{
			Password: password,
		},
	}
	newJob := models.NewJob{
		JobName:     "Set password",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opSetPassword.ToResource(),
		},
	}
	return newJob
}

func TaskRemovePassword(password []byte) models.NewJob {
	opAuth := models.JobStep{
		Command: models.CommandAuthPassword,
		Params: models.AuthPasswordParams{
			Password: password,
		},
	}

	opRemovePassword := models.JobStep{
		Command: models.CommandRemovePassword,
	}

	newJob := models.NewJob{
		JobName:     "Remove password",
		Repeat:      1,
		ExpireAfter: 60,
		Steps: []models.JobStepResource{
			opAuth.ToResource(),
			opRemovePassword.ToResource(),
		},
	}
	return newJob
}
