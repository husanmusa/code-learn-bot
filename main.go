package main

import (
	"context"
	"flag"
	"github.com/husanmusa/code-learn-bot/bot"
	"github.com/husanmusa/code-learn-bot/config"
	"github.com/husanmusa/code-learn-bot/pkg/db"
	"github.com/husanmusa/code-learn-bot/pkg/logger"
	"github.com/husanmusa/code-learn-bot/service/lesson"
	"github.com/husanmusa/code-learn-bot/service/user"
	lessonStorage "github.com/husanmusa/code-learn-bot/storage/postgres/lesson"
	userStorage "github.com/husanmusa/code-learn-bot/storage/postgres/user"

	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

func init() {
	err := os.Setenv("TZ", "Asia/Tashkent")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logs := logger.New("develop", "CODING")
	err := godotenv.Load(".env")
	if err != nil {
		logs.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			logs.Fatal("failed cleanup logger", logger.Error(err))
		}
	}(logs)

	err = setupLogging(logs)
	if err != nil {
		logs.Fatal("failed setup Logging", logger.Error(err))
	}
	logs.Info("main: sqlxConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase),
	)

	conDB, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sqlx connection to postgres error", logger.Error(err))
	}

	defer conDB.Close()

	userService := user.NewService(userStorage.New(conDB))
	lessonService := lesson.NewService(lessonStorage.New(conDB))

	go func() {
		err = bot.Start(
			context.Background(),
			cfg.BotToken,
			userService,
			lessonService,
		)
		if err != nil {
			log.Println(err)
		}
	}()
	log.Printf("start listening and serving at %q\n", "8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func setupLogging(cfg *logger.LoggerImpl) error {

	var (
		logPath   = flag.String("log", "./logs.txt", "path to log file")
		logOutput io.Writer
	)

	logOutput, err := os.OpenFile(*logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(logOutput)
	return nil
}
