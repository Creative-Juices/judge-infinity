package services

import (
	"io"
	"judgeinf/internal/models"

	"github.com/google/uuid"
)

type Queue interface {
	PushRequestToQueue(SubmissionID uuid.UUID) error
}

type Table interface {
	PushRequestMetadataToTable(CodeSubmission models.CodeSubmission) error
	FetchRequestMetadataFromTable(SubmissionID uuid.UUID) (models.CodeSubmission, error)
	PushQuestionMetadataToTable(Question models.Question) error
	FetchQuestionMetadataToTable(QuestionID uuid.UUID) (models.Question, error)
	PushResultsToTable(SubmissionResult models.SubmissionResult) error
	FetchResultsFromTable(SubmissionID uuid.UUID) (models.SubmissionResult, error)
}

type Storage interface {
	PushTestcaseToStorage(Testcase io.Reader) error
}
