package main

import (
	"archive/zip"
	"errors"
	"judgeinf/internal/models"
	"judgeinf/internal/services"
	"os"
	"strconv"
	"strings"
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

	InitServices()
	InitLogger()

	bindRoutes(api)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen("127.0.0.1:3000")
}

func bindRoutes(r fiber.Router) {
	r.Post("/questions/createOrUpdate", createOrUpdateQuestion)
	r.Get("/submissions/:submissionId/status", getSubmissionStatus)
	r.Post("/submissions/submit", makeSubmission)
}

func createOrUpdateQuestion(c *fiber.Ctx) error {
	question := &models.Question{}
	var questionId uuid.UUID
	var timeLimitMultiplier int
	var err error

	questionIdStr := strings.TrimSpace(c.FormValue("question_id"))
	timeLimitMultiplierStr := strings.TrimSpace(c.FormValue("time_limit_multiplier"))

	if len(questionIdStr) > 0 {
		questionId, err = uuid.Parse(questionIdStr)
		if err != nil {
			c.SendString("Unexpected error")
			return err
		}
	} else {
		questionId = uuid.New()
	}
	question.QuestionID = questionId

	if len(timeLimitMultiplierStr) > 0 {
		timeLimitMultiplier, err = strconv.Atoi(timeLimitMultiplierStr)
		if err != nil {
			c.SendString("Unexpected error")
			return err
		}
	} else {
		timeLimitMultiplier = 1
	}
	question.TimeLimitMultiplier = timeLimitMultiplier

	testcasesZip, err := c.FormFile("testcases")
	if err != nil {
		c.SendString("Unexpected error")
		return err
	}

	testcasesZipFile, err := testcasesZip.Open()
	if err != nil {
		c.SendString("Unexpected error")
		return err
	}

	testcasesZipFileReader, err := zip.NewReader(testcasesZipFile, testcasesZip.Size)
	if err != nil {
		c.SendString("Unexpected error")
		return err
	}

	testcaseFiles, err := services.StorageInstance.PushTestcasesToStorage(testcasesZipFileReader.File, question.QuestionID)
	if err != nil {
		c.SendString("Unexpected error")
		return err
	}
	question.Testcases = testcaseFiles

	log.Logger.Print(question)

	err = services.TableInstance.PushQuestionMetadataToTable(question)
	if err != nil {
		c.SendString("Unexpected error")
		return err
	}

	if err := c.JSON(question); err != nil {
		c.SendString("Unexpected error")
		return err
	}

	return nil
}

func getSubmissionStatus(c *fiber.Ctx) error {
	submissionIdStr := c.Params("submissionId")
	if submissionIdStr == "" {
		c.SendString("Invalid request")
		return errors.New("invalid request")
	}

	submissionId, err := uuid.Parse(submissionIdStr)
	if err != nil {
		c.SendString("Invalid request")
		return err
	}

	submission, err := services.TableInstance.FetchResultsFromTable(submissionId)
	if err != nil {
		c.SendString("Invalid request")
		return err
	}

	if err := c.JSON(submission); err != nil {
		c.SendString("Invalid request")
		return err
	}

	return nil
}

func makeSubmission(c *fiber.Ctx) error {
	submission := &models.CodeSubmission{}
	response := &models.StandardResponse{}

	if err := c.BodyParser(submission); err != nil {
		c.SendString("Unexpected error")
		return err
	}

	submission.SubmissionID = uuid.New()
	if err := services.TableInstance.PushRequestMetadataToTable(submission); err != nil {
		c.SendString("Unexpected error")
		return err
	}
	if err := services.QueueInstance.PushRequestToQueue(submission.SubmissionID); err != nil {
		c.SendString("Unexpected error")
		return err
	}

	response.Success = true
	response.Message = "done"
	response.Value = submission.SubmissionID
	if err := c.JSON(response); err != nil {
		c.SendString("Unexpected error")
		return err
	}

	return nil
}

const (
	requestIdKey = "requestId"
)

func InitServices() {
	queue, err := services.NewQueue()
	if err != nil {
		panic(err)
	}

	table, err := services.NewTable()
	if err != nil {
		panic(err)
	}

	storage, err := services.NewStorage()
	if err != nil {
		panic(err)
	}

	languageData, err := services.NewLanguageData()
	if err != nil {
		panic(err)
	}

	services.QueueInstance = queue
	services.TableInstance = table
	services.StorageInstance = storage
	services.LanguageData = languageData
}

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
