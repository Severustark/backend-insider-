package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func Init() {
	// Log seviyesini ayarla (opsiyonel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Global logger'ı ayarlamak istersen:
	log.Logger = Log
}
