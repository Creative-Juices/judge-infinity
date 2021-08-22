package services

import (
	"io"
	"judgeinf/internal/models"

	"github.com/google/uuid"
)

type Queue interface {
	PushRequestToQueue(SubmissionID uuid.UUID)
}

type Table interface {
	PushRequestMetadataToTable(CodeSubmission models.CodeSubmission)
	FetchRequestMetadataFromTable(SubmissionID uuid.UUID) models.CodeSubmission
	PushQuestionMetadataToTable(Question models.Question)
	FetchQuestionMetadataToTable(QuestionID uuid.UUID) models.Question
	PushResultsToTable(SubmissionResult models.SubmissionResult)
	FetchResultsFromTable(SubmissionID uuid.UUID) models.SubmissionResult
}

type Storage interface {
	PushTestcaseToStorage(Testcase io.Reader)
}
