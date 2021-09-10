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
var entityFetchTimeout = uint(30)

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
	return azureQueue.SubmissionQueue.GetMessageReference(SubmissionID.String()).Put(&storage.PutMessageOptions{})
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

func (azureTable *AzureTable) PushRequestMetadataToTable(CodeSubmission *models.CodeSubmission) error {
	submissionIdStr := CodeSubmission.SubmissionID.String()
	languageIdStr := CodeSubmission.LanguageID.String()
	questionIdStr := CodeSubmission.QuestionID.String()

	entity := azureTable.RequestMetadata.GetEntityReference(submissionIdStr, submissionIdStr)

	props := map[string]interface{}{
		"SourceCode":   CodeSubmission.SourceCode,
		"LanguageID":   languageIdStr,
		"QuestionID":   questionIdStr,
		"CallbackUrl":  CodeSubmission.CallbackUrl,
		"ResponseMode": CodeSubmission.ResponseMode,
		"SubmissionID": submissionIdStr,
	}
	entity.Properties = props

	return entity.InsertOrReplace(nil)
}

func (azureTable *AzureTable) FetchRequestMetadataFromTable(SubmissionID uuid.UUID) (*models.CodeSubmission, error) {
	submission := &models.CodeSubmission{}
	submissionIdStr := SubmissionID.String()

	entity := azureTable.RequestMetadata.GetEntityReference(submissionIdStr, submissionIdStr)

	if err := entity.Get(entityFetchTimeout, storage.NoMetadata, nil); err != nil {
		return nil, err
	}

	entityAsJson, err := json.Marshal(entity.Properties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(entityAsJson, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

func (azureTable *AzureTable) PushQuestionMetadataToTable(Question *models.Question) error {
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

	return entity.InsertOrReplace(nil)
}

func (azureTable *AzureTable) FetchQuestionMetadataFromTable(QuestionID uuid.UUID) (*models.Question, error) {
	questionIdStr := QuestionID.String()
	question := &models.Question{}

	entity := azureTable.QuestionMetadata.GetEntityReference(questionIdStr, questionIdStr)

	if err := entity.Get(entityFetchTimeout, storage.NoMetadata, nil); err != nil {
		return nil, err
	}

	/*if value, ok := entity.Properties["QuestionID"]; !ok {
		return nil, errors.New("question not found in table")
	} else if questionId, err := uuid.Parse(fmt.Sprintf("%v", value)); err != nil {
		return nil, err
	} else {
		question.QuestionID = questionId
	}

	if value, ok := entity.Properties["TimeLimitMultiplier"]; !ok {
		return nil, errors.New("question not found in table")
	} else if timeLimitMultiplier, err := strconv.Atoi(fmt.Sprintf("%v", value)); err != nil {
		return nil, err
	} else {
		question.TimeLimitMultiplier = timeLimitMultiplier
	}

	if value, ok := entity.Properties["Testcases"]; !ok {
		return nil, errors.New("question not found in table")
	} else if err := json.Unmarshal([]byte(fmt.Sprintf("%v", value)), &(question.Testcases)); err != nil {
		return nil, err
	}*/

	entityAsJson, err := json.Marshal(entity.Properties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(entityAsJson, question); err != nil {
		return nil, err
	}

	fmt.Println("Parsed question from table")
	fmt.Println(question)

	return question, nil
}

func (azureTable *AzureTable) PushResultsToTable(SubmissionResult *models.SubmissionResult) error {
	submissionIdStr := SubmissionResult.SubmissionID.String()

	entity := azureTable.Results.GetEntityReference(submissionIdStr, submissionIdStr)

	props := map[string]interface{}{
		"SubmissionID": submissionIdStr,
		"Verdicts":     SubmissionResult.Verdicts,
	}
	entity.Properties = props

	return entity.InsertOrReplace(nil)
}

func (azureTable *AzureTable) FetchResultsFromTable(SubmissionID uuid.UUID) (*models.SubmissionResult, error) {
	submissionIdStr := SubmissionID.String()
	submissionResult := &models.SubmissionResult{}

	entity := azureTable.Results.GetEntityReference(submissionIdStr, submissionIdStr)

	if err := entity.Get(entityFetchTimeout, storage.NoMetadata, nil); err != nil {
		return nil, err
	}

	entityAsJson, err := json.Marshal(entity.Properties)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(entityAsJson, submissionResult); err != nil {
		return nil, err
	}

	return submissionResult, nil
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
	testcases := []string{}
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
					return nil, err
				}
				testcaseOutputFileReader, err := file.Open()
				if err != nil {
					return nil, err
				}

				err = testcaseInputFileBlobRef.CreateBlockBlobFromReader(testcaseInputFileReader, nil)
				if err != nil {
					return nil, err
				}
				err = testcaseOutputFileBlobRef.CreateBlockBlobFromReader(testcaseOutputFileReader, nil)
				if err != nil {
					return nil, err
				}

				testcases = append(testcases, fileName)
			}
		}
	}

	return testcases, nil
}
