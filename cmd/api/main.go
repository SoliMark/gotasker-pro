package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SoliMark/gotasker-pro/config"
)

func main() {
	// 1. 載入設定
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("faile to load config: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := fmt.Fprintln(w, "Hello from GoTasker Pro!"); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})

	fmt.Printf("Server running on port %s\n", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, nil))
}
