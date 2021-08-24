package models

import (
	"github.com/google/uuid"
)

type Question struct {
	QuestionID          uuid.UUID `json:"question_id"`
	TimeLimitMultiplier int       `json:"time_limit_multiplier"`
	Testcases           []string  `json:"testcases,omitempty"`
}

type ResponseModeType int

const (
	AfterEachTestcase ResponseModeType = 0
	AfterAllTestcases ResponseModeType = 1
)

type CodeSubmission struct {
	SourceCode   string           `json:"source_code"`
	LanguageID   uuid.UUID        `json:"language_id"`
	QuestionID   uuid.UUID        `json:"question_id"`
	CallbackUrl  string           `json:"callback_url"`
	ResponseMode ResponseModeType `json:"response_mode_type"`
	SubmissionID uuid.UUID        `json:"submission_id,omitempty"`
}

type GetResultsRequest struct {
	SubmissionID uuid.UUID `json:"submission_id"`
}

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Value   interface{} `json:"value"`
}

type TestcaseVerdict string

const (
	CorrectAnswer       TestcaseVerdict = "AC"
	WrongAnswer         TestcaseVerdict = "WA"
	RuntimeError        TestcaseVerdict = "RE"
	Processing          TestcaseVerdict = "PROC"
	TimeLimitExceeded   TestcaseVerdict = "TLE"
	MemoryLimitExceeded TestcaseVerdict = "MLE"
	CompilationError    TestcaseVerdict = "CE"
)

type SubmissionResult struct {
	SubmissionID uuid.UUID                  `json:"submission_id"`
	Verdicts     map[string]TestcaseVerdict `json:"verdicts"`
}
