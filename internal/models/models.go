package models

import (
	"github.com/google/uuid"
)

type Question struct {
	QuestionID          uuid.UUID
	TimeLimitMultiplier int
	Testcases           []string
}

type ResponseModeType int

const (
	AfterEachTestcase ResponseModeType = 0
	AfterAllTestcases ResponseModeType = 1
)

type CodeSubmission struct {
	SourceCode   string
	LanguageID   uuid.UUID
	QuestionID   uuid.UUID
	CallbackUrl  string
	ResponseMode ResponseModeType
	SubmissionID uuid.UUID
}

type GetResultsRequest struct {
	SubmissionID uuid.UUID
}

type StandardResponse struct {
	Success bool
	Message string
	Value   interface{}
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
	SubmissionID uuid.UUID
	Verdicts     map[string]TestcaseVerdict
}
