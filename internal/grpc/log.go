package grpcsrv

import (
	"io"
	"log"
	"os"
)

var AccessLogger *log.Logger

func SetupAccessLogger() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(
		"logs/access.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	mw := io.MultiWriter(os.Stdout, file)

	AccessLogger = log.New(
		mw,
		"",
		log.LstdFlags,
	)

	return nil
}


