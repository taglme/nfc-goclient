package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/taglme/nfc-goclient/pkg/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	baseURL := strings.TrimSpace(os.Getenv("BASE_URL"))
	if baseURL == "" {
		baseURL = "http://127.0.0.1:3011"
	}
	xAppKey := strings.TrimSpace(os.Getenv("X_APP_KEY"))
	xUserToken := strings.TrimSpace(os.Getenv("X_USER_TOKEN"))
	locale := strings.TrimSpace(os.Getenv("LOCALE"))

	c := client.NewWithOptions(client.Options{
		BaseURL:      baseURL,
		Locale:       locale,
		XAppKey:      xAppKey,
		XUserToken:   xUserToken,
		Interceptors: nil,
	})
	adapters, err := c.Adapters.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(adapters)
}
