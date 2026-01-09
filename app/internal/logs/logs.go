package logs

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var logger *log.Logger

func init() {
	cwd, _ := os.Getwd()
	logDir := filepath.Join(cwd, "logs")
	os.MkdirAll(logDir, 0755)

	logPath := filepath.Join(logDir, "healthchecker.log")
	file, err := os.OpenFile(
		logPath,
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
