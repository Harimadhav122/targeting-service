package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting Application..")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error("Failed to start server: ", err)
	}
}
