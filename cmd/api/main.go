package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	app := fiber.New()
	app.Use(RequestLogger)

	api := app.Group("/api/v1")

	InitLogger()

	bindRoutes(api)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen("127.0.0.1:3000")
}

func bindRoutes(r fiber.Router) {
	r.Post("/questions/create", createQuestion)
	r.Get("/submissions/:submissionId/status", getSubmissionStatus)
	r.Post("/submissions/submit/:questionId", makeSubmission)
}

func createQuestion(c *fiber.Ctx) error {
	panic("Unimplemented")
}

func getSubmissionStatus(c *fiber.Ctx) error {
	panic("Unimplemented")

}

func makeSubmission(c *fiber.Ctx) error {
	panic("Unimplemented")
}

const (
	requestIdKey = "requestId"
)

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func RequestLogger(c *fiber.Ctx) error {

	start := time.Now()

	id := uuid.New()
	logger := log.With().Str(requestIdKey, id.String()).Logger()
	c.Locals(requestIdKey, &logger)
	msg := "Request"
	if err := c.Next(); err != nil {
		msg = err.Error()
		log.Error().Stack().Str(requestIdKey, id.String()).Err(err).Msg("")
	}

	code := c.Response().StatusCode()

	apiLogger := log.With().
		Str("method", c.Method()).
		Str("path", c.Path()).
		Str("ip", c.IP()).
		Str("responseTime", time.Since(start).String()).
		Int("status", code).Logger()

	switch {
	case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
		apiLogger.Warn().Msg(msg)
	case code >= fiber.StatusInternalServerError:
		apiLogger.Error().Msg(msg)
	default:
		apiLogger.Info().Msg(msg)
	}

	return nil
}
