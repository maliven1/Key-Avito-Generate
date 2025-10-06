package main

import (
	"avito/internal/config"
	"avito/internal/handler"
	"avito/internal/logger"
	"avito/internal/logic"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.NewConfig()

	str := slog.String("env", cfg.Env)

	logger.SetupLogger(cfg.Env)
	logger.Log.Info("Starting project", str)
	logger.Log.Debug("debug messages are enabled", str)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wg := sync.WaitGroup{}
	counter := 3
	wg.Add(counter)
	for i := 0; i < counter; i++ {
		titleKey, keys := checkNewKeys("электрик")

		go func() {
			handler.JptHandler(keys, titleKey, i)
			wg.Done()
		}()
	}
	wg.Wait()
}

func checkNewKeys(key string) (string, string) {
	switch key {
	case "сантехник":
		return logic.Logic(os.Getenv("SANTEX_STORAGE"), os.Getenv("SANTEX_PATH"))
	case "электрик":
		return logic.Logic(os.Getenv("ELECTR_STORAGE"), os.Getenv("ELECTR_PATH"))
	default:
		logger.Log.Error("Invalid key")
		return "", ""
	}
}
