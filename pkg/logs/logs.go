package logs

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func Init(env string) {
	if env == "test" {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		Logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
		return
	}

	if _, err := os.Stat("storage"); os.IsNotExist(err) {
		err := os.Mkdir("storage", 0755)
		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("storage/logs"); os.IsNotExist(err) {
		err := os.Mkdir("storage/logs", 0755)
		if err != nil {
			panic(err)
		}
	}

	file, err := os.OpenFile("storage/logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	Logger = slog.New(slog.NewTextHandler(file, opts))
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}
