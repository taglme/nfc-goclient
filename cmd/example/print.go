package main

import (
	"fmt"

	"github.com/taglme/nfc-goclient/pkg/models"
)

func printInfo(info models.AppInfo) {
	fmt.Printf("Get server info: version - %s, platform - %s\n", info.Version, info.Platform)
}
func printAdapters(adapters []models.Adapter) {
	if len(adapters) == 0 {
		fmt.Println("Adapters not available")
		return
	}
	fmt.Println("Current adapters:")
	for i, adapter := range adapters {
		fmt.Printf("[%d] %s\n", i, adapter.Name)
	}
	return
}
