package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Service interface {
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(err error, fields ...map[string]interface{})
	Fatal(msg string, fields ...map[string]interface{})
	Panic(msg string, fields ...map[string]interface{})
}

type service struct {
	logger *zerolog.Logger
}

func New() Service {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	)

	return &service{&log}
}

func (s *service) Debug(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Debug().Fields(f).Msg(msg)
		return
	}

	s.logger.Debug().Msg(msg)
}

func (s *service) Info(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Info().Fields(f).Msg(msg)
		return
	}

	s.logger.Info().Msg(msg)
}

func (s *service) Warn(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Warn().Fields(f).Msg(msg)
		return
	}

	s.logger.Warn().Msg(msg)
}

func (s *service) Error(err error, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Error().Stack().Err(err).Fields(f).Msg("")
		return
	}

	s.logger.Error().Stack().Err(err).Msg("")
}

func (s *service) Fatal(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Fatal().Fields(f).Msg(msg)
		return
	}

	s.logger.Fatal().Msg(msg)
}

func (s *service) Panic(msg string, fields ...map[string]interface{}) {
	if len(fields) > 0 {
		f := s.mergeFields(fields...)
		s.logger.Panic().Fields(f).Msg(msg)
		return
	}

	s.logger.Panic().Msg(msg)
}

func (s *service) mergeFields(fields ...map[string]interface{}) map[string]interface{} {
	f := make(map[string]interface{})
	for _, field := range fields {
		for key, value := range field {
			f[key] = value
		}
	}
	return f
}
