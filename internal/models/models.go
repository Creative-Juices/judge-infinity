package models

import (
	"github.com/google/uuid"
)

type Question struct {
	QuestionID          uuid.UUID `json:"question_id,omitempty"`
	TimeLimitMultiplier int       `json:"time_limit_multiplier,omitempty"`
	Testcases           []string  `json:"testcases,omitempty"`
}

type ResponseModeType int

const (
	AfterEachTestcase ResponseModeType = 0
	AfterAllTestcases ResponseModeType = 1
)

type Language string

const (
	Python Language = "python"
	Java   Language = "java"
	CPP    Language = "cpp"
)

type CodeSubmission struct {
	SourceCode   string           `json:"source_code"`
	Language     Language         `json:"language"`
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

type LanguageInfo struct {
	TargetFilename string `json:"target_filename"`
}
