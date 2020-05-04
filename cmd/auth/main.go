package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/taglme/nfc-goclient/pkg/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	privateRSAKey, err := client.PrivateRSAKeyFromB64String(os.Getenv("SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	auth := client.NewSigner(os.Getenv("APP_ID"), privateRSAKey, os.Getenv("CERT"))
	client := client.New("127.0.0.1:3011", auth)
	adapters, err := client.Adapters.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(adapters)
}
