package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Caller().Logger()

	// switch strings.ToLower(config.LogLevel) {
	// case "trace":
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	// case "debug":
	// 	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// case "info":
	// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// case "warn":
	// 	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	// case "error":
	// 	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	// case "fatal":
	// 	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	// case "panic":
	// 	zerolog.SetGlobalLevel(zerolog.PanicLevel)
	// default:
	// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// }
}
