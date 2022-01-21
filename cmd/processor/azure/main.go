package main

import (
	"encoding/json"
	"fmt"
	"judgeinf/internal/models"
	"judgeinf/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var basePath = "/home/site/wwwroot"

type InvokeRequest struct {
	Data     map[string]string
	Metadata map[string]interface{}
}

func IncomingQueueTrigger(w http.ResponseWriter, r *http.Request) {
	invokeRequest := InvokeRequest{}
	submissionId := uuid.UUID{}
	response := models.StandardResponse{}

	if err := json.NewDecoder(r.Body).Decode(&invokeRequest); err != nil {
		fmt.Println(err)
		return
	}

	item := invokeRequest.Data["myQueueItem"]
	if err := json.Unmarshal([]byte(item), &submissionId); err != nil {
		fmt.Println(err)
		return
	}

	requestMetadata, err := services.TableInstance.FetchRequestMetadataFromTable(submissionId)
	if err != nil {
		fmt.Println(err)
		return
	}
	questionMetadata, err := services.TableInstance.FetchQuestionMetadataFromTable(requestMetadata.QuestionID)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(submissionId)
	fmt.Println(questionMetadata)
	fmt.Println(requestMetadata)

	// Initialize verdicts
	submissionResult := &models.SubmissionResult{
		SubmissionID: requestMetadata.SubmissionID,
		Verdicts:     make(map[string]models.TestcaseVerdict),
	}
	for _, testcase := range questionMetadata.Testcases {
		submissionResult.Verdicts[testcase] = models.RuntimeError
	}
	if err := services.TableInstance.PushResultsToTable(submissionResult); err != nil {
		fmt.Println(err)
		return
	}

	languageInfo := services.LanguageData[requestMetadata.Language]

	// Create source file
	if err := os.WriteFile(fmt.Sprintf("%s/%s", basePath, languageInfo.TargetFilename), []byte(requestMetadata.SourceCode), 0777); err != nil {
		fmt.Println(err)
		return
	}

	response.Success = true
	response.Message = "done"
	response.Value = submissionId
	responseAsJson, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseAsJson)
}

func main() {
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	InitServices()

	r := http.NewServeMux()
	r.HandleFunc("/IncomingQueueTrigger", IncomingQueueTrigger)
	fmt.Println("Listening on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func InitServices() {
	queue, err := services.NewQueue()
	if err != nil {
		panic(err)
	}

	table, err := services.NewTable()
	if err != nil {
		panic(err)
	}

	storage, err := services.NewStorage()
	if err != nil {
		panic(err)
	}

	languageData, err := services.NewLanguageData()
	if err != nil {
		panic(err)
	}

	services.QueueInstance = queue
	services.TableInstance = table
	services.StorageInstance = storage
	services.LanguageData = languageData
}
