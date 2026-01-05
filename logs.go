package main

import (
	"io"
	"log"
	"os"
)

var logger *log.Logger

func init() {
	file, err := os.OpenFile(
		"healthcheck.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	logger = log.New(
		mw,
		"",
		log.LstdFlags,
	)
}

func LogEvent(
	message string,
) {
	logger.Println(message)
}
