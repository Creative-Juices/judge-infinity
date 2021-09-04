package azure

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"judgeinf/internal/models"
	"strings"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/google/uuid"
)

var storageConnectionString = ""
var submissionQueueName = "submissionqueue"
var requestMetadataTable = "requestmetadata"
var questionMetadataTable = "questionmetadata"
var resultsTable = "results"
var testcaseStorageContainer = "testcases"
var inputPrefix = "input/"
var outputPrefix = "output/"

type AzureQueue struct {
	SubmissionQueue *storage.Queue
}

func NewAzureQueue() (*AzureQueue, error) {
	client, err := storage.NewClientFromConnectionString(storageConnectionString)

	if err != nil {
		panic(err)
	}

	queueService := client.GetQueueService()

	azureQueue := &AzureQueue{
		SubmissionQueue: queueService.GetQueueReference(submissionQueueName),
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
	client, err := storage.NewClientFromConnectionString(storageConnectionString)

	if err != nil {
		panic(err)
	}

	tableService := client.GetTableService()

	azureTable := &AzureTable{
		RequestMetadata:  tableService.GetTableReference(requestMetadataTable),
		QuestionMetadata: tableService.GetTableReference(questionMetadataTable),
		Results:          tableService.GetTableReference(resultsTable),
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
	questionIdStr := Question.QuestionID.String()

	entity := azureTable.QuestionMetadata.GetEntityReference(questionIdStr, questionIdStr)

	testcasesJsonMarshaled, err := json.Marshal(Question.Testcases)
	if err != nil {
		return err
	}

	props := map[string]interface{}{
		"QuestionID":          questionIdStr,
		"TimeLimitMultiplier": Question.TimeLimitMultiplier,
		"Testcases":           string(testcasesJsonMarshaled),
	}
	entity.Properties = props

	err = entity.InsertOrReplace(&storage.EntityOptions{})
	return err
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
	client, err := storage.NewClientFromConnectionString(storageConnectionString)

	if err != nil {
		panic(err)
	}

	blobService := client.GetBlobService()

	azureStorage := &AzureStorage{
		TestcaseStorage: blobService.GetContainerReference(testcaseStorageContainer),
	}

	return azureStorage, nil
}

func (azureStorage *AzureStorage) PushTestcasesToStorage(Testcases []*zip.File, QuestionID uuid.UUID) ([]string, error) {
	inputFilesFound := make(map[string]*zip.File)
	var testcases []string
	questionIdStr := QuestionID.String()

	for _, file := range Testcases {
		fileInfo := file.FileInfo()
		fileName := fileInfo.Name()
		if !fileInfo.IsDir() && strings.HasPrefix(file.Name, inputPrefix) {
			inputFilesFound[fileName] = file
		}
	}

	for _, file := range Testcases {
		fileInfo := file.FileInfo()
		fileName := fileInfo.Name()
		if !fileInfo.IsDir() && strings.HasPrefix(file.Name, outputPrefix) {
			if inputFile, ok := inputFilesFound[fileName]; ok {
				testcaseInputFileName := fmt.Sprintf("%s/%s%s", questionIdStr, inputPrefix, fileName)
				testcaseOutputFileName := fmt.Sprintf("%s/%s%s", questionIdStr, outputPrefix, fileName)

				testcaseInputFileBlobRef := azureStorage.TestcaseStorage.GetBlobReference(testcaseInputFileName)
				testcaseOutputFileBlobRef := azureStorage.TestcaseStorage.GetBlobReference(testcaseOutputFileName)

				testcaseInputFileReader, err := inputFile.Open()
				if err != nil {
					return []string{}, err
				}
				testcaseOutputFileReader, err := file.Open()
				if err != nil {
					return []string{}, err
				}

				err = testcaseInputFileBlobRef.CreateBlockBlobFromReader(testcaseInputFileReader, &storage.PutBlobOptions{})
				if err != nil {
					return []string{}, err
				}
				err = testcaseOutputFileBlobRef.CreateBlockBlobFromReader(testcaseOutputFileReader, &storage.PutBlobOptions{})
				if err != nil {
					return []string{}, err
				}

				testcases = append(testcases, fileName)
			}
		}
	}

	return testcases, nil
}
