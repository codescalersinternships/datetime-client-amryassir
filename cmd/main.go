package main

import (
	"fmt"
	"log"

	httpclient "github.com/codescalersinternships/datetime-client-amryassir/pkg"
)

func main() {
	config := httpclient.LoadConfig()
	httpClient := httpclient.NewClient(config)

	operation := func() error {
		dateTime, err := httpClient.GetDateTime()
		if err != nil {
			return err
		}
		fmt.Println("Current DateTime:", dateTime)
		return nil
	}

	if err := httpclient.Retry(operation); err != nil {
		log.Fatalf("Failed to get datetime: %v", err)
	}
}
