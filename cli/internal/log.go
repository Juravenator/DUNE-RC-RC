package internal

import (
	"os"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

var log = logger.With().Str("pkg", "internal").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
