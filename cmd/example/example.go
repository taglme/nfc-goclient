package main

import (
	"fmt"

	"github.com/taglme/nfc-client/pkg/client"
)

func main() {
	client := client.New("http://127.0.0.1:3011")
	adapters, err := client.Adapters.GetAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	about, err := client.About.Get()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(about.Name)
	for i, adapter := range adapters {
		adapterWithDetails, err := client.Adapters.Get(adapter.AdapterID)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Adapter %d: %s, driver: %s\n", i, adapterWithDetails.Name, adapterWithDetails.Driver)
	}
}
