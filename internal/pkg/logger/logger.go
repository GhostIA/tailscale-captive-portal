package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

type LogAdapter struct {
	logger *slog.Logger
}

func NewLogger() *LogAdapter {
	return &LogAdapter{
		logger: slog.New(
			tint.NewHandler(os.Stdout, &tint.Options{
				Level:      slog.LevelDebug,
				TimeFormat: time.Kitchen,
				AddSource:  true,
			}),
		),
	}
}

func (log *LogAdapter) ReturnDefault() *slog.Logger {
	return log.logger
}
func (log *LogAdapter) Printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)

	if strings.Contains(strings.ToLower(msg), "error") {
		log.logger.Error(msg)
	} else {
		log.logger.Info(msg)
	}
}
