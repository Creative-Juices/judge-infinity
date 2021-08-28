package azure

import (
	"io"
	"judgeinf/internal/models"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/google/uuid"
)

type AzureQueue struct {
	SubmissionQueue *storage.Queue
}

func NewAzureQueue() (*AzureQueue, error) {
	client, err := storage.NewClientFromConnectionString("")

	if err != nil {
		panic(err)
	}

	queueService := client.GetQueueService()

	azureQueue := &AzureQueue{
		SubmissionQueue: queueService.GetQueueReference(""),
	}

	return azureQueue, nil
}

func (azureQueue *AzureQueue) PushRequestToQueue(SubmissionID uuid.UUID) error {
	return nil
}

type AzureTable struct {
	RequestMetadata  *storage.Table
	QuestionMetadata *storage.Table
	Results          *storage.Table
}

func NewAzureTable() (*AzureTable, error) {
	client, err := storage.NewClientFromConnectionString("")

	if err != nil {
		panic(err)
	}

	tableService := client.GetTableService()

	azureTable := &AzureTable{
		RequestMetadata:  tableService.GetTableReference(""),
		QuestionMetadata: tableService.GetTableReference(""),
		Results:          tableService.GetTableReference(""),
	}

	return azureTable, nil
}

func (azureTable *AzureTable) PushRequestMetadataToTable(CodeSubmission models.CodeSubmission) error {
	return nil
}

func (azureTable *AzureTable) FetchRequestMetadataFromTable(SubmissionID uuid.UUID) (models.CodeSubmission, error) {
	return models.CodeSubmission{}, nil
}

func (azureTable *AzureTable) PushQuestionMetadataToTable(Question models.Question) error {
	return nil
}

func (azureTable *AzureTable) FetchQuestionMetadataToTable(QuestionID uuid.UUID) (models.Question, error) {
	return models.Question{}, nil
}

func (azureTable *AzureTable) PushResultsToTable(SubmissionResult models.SubmissionResult) error {
	return nil
}

func (azureTable *AzureTable) FetchResultsFromTable(SubmissionID uuid.UUID) (models.SubmissionResult, error) {
	return models.SubmissionResult{}, nil
}

type AzureStorage struct {
	TestcaseStorage *storage.Container
}

func NewAzureStorage() (*AzureStorage, error) {
	client, err := storage.NewClientFromConnectionString("")

	if err != nil {
		panic(err)
	}

	blobService := client.GetBlobService()

	azureStorage := &AzureStorage{
		TestcaseStorage: blobService.GetContainerReference(""),
	}

	return azureStorage, nil
}

func (azureStorage *AzureStorage) PushTestcaseToStorage(Testcase io.Reader) error {
	return nil
}
