package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	baseURL := strings.TrimSpace(os.Getenv("BASE_URL"))
	if baseURL == "" {
		baseURL = "http://127.0.0.1:3011"
	}
	xAppKey := strings.TrimSpace(os.Getenv("X_APP_KEY"))
	xUserToken := strings.TrimSpace(os.Getenv("X_USER_TOKEN"))
	locale := strings.TrimSpace(os.Getenv("LOCALE"))

	tasker := NewTasker(baseURL, xAppKey, xUserToken, locale)
	tasker.Add(TaskWriteURL("https://tagl.me"))
	tasker.Add(TaskRead())
	tasker.Add(TaskGetDump())
	tasker.Add(TaskTransmit([]byte{0x60}))
	tasker.Add(TaskFormatDefault())
	tasker.Add(TaskSetPassword([]byte{0x11, 0x11, 0x11, 0x11}))
	tasker.Add(TaskRemovePassword([]byte{0x11, 0x11, 0x11, 0x11}))

	err := tasker.Run(1)
	if err != nil {
		log.Printf("Tasks execution failed with error: %s", err)
	}
	return

}
