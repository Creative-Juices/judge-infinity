package services

import (
	"archive/zip"
	"judgeinf/internal/models"
	"judgeinf/internal/services/azure"

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
	PushTestcasesToStorage(Testcases []*zip.File, QuestionID uuid.UUID) ([]string, error)
}

func NewQueue() (Queue, error) {
	return azure.NewAzureQueue()
}

func NewTable() (Table, error) {
	return azure.NewAzureTable()
}

func NewStorage() (Storage, error) {
	return azure.NewAzureStorage()
}

var (
	QueueInstance   Queue
	TableInstance   Table
	StorageInstance Storage
)
