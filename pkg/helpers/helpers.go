package helpers

import (
	log "github.com/sirupsen/logrus"
)

// ErrorHandler  Display message & error
func ErrorHandler(details string, err error) {
	if err != nil {
		log.Println(details, err)
	}
}

// ErrorHandlerFatal Display message & error & terminate
func ErrorHandlerFatal(details string, err error) {
	if err != nil {
		log.Fatal(details, err)
	}
}
