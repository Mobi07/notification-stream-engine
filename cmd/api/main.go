package main

import "github.com/Mobi07/notification-stream-engine.git/pkg/logger"

func main() {
	logger.Init()
	defer logger.Sync()

	logger.Log.Info("API service started")
}