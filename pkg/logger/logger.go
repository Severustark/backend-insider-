package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger

func Init() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()

	log.Logger = Log
}
