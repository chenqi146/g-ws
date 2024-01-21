package main

import (
	"go.uber.org/zap"
	"log"
)

var l *zap.Logger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	l = logger
}
