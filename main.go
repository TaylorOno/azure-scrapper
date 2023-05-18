package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	fmt.Println("scrapping azure")
	sub := os.Getenv("AZURE_SUBSCRIPTION")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}

	scrapper, err := NewScrapper(cred, sub)
	if err != nil {
		log.Fatalf("unable to initialize scrapper: %v", err)
	}

	if err = scrapper.Run(); err != nil {
		log.Fatalf("scrapper failed: %v", err)
	}
}

func consoleHandler[T any](t *T) error {
	return json.NewEncoder(os.Stdout).Encode(t)
}
