package daq

import (
	"os"

	"github.com/rs/zerolog"
	logger "github.com/rs/zerolog/log"
)

var log = logger.With().Str("pkg", "daq").Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
