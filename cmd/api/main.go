package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	api := app.Group("/api/v1")

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
