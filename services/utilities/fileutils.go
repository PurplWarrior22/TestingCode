package utilities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/PurplWarrior22/TestingCode/services/models"
)

//general convenience utilities for working with files

const (
	//TODO: make these errors a little more generic than our data service error
	FILE_READ_ERROR    = "Error reading feed data"
	FILE_OPEN_ERROR    = "Error finding feed data"
	FEED_DID_NOT_EXIST = "Feed with the passed id did not exist"
	UNMARSHALL_ERROR   = "Unable to marshall data into the passed object"
)

func ReadBytesOrError(filePath string) ([]byte, error) {
	feedsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error while opening feed file %s\n", err.Error())
		if os.IsNotExist(err) {
			return nil, models.NewDataServiceError(FEED_DID_NOT_EXIST, http.StatusNotFound)
		} else {
			return nil, models.NewDataServiceError(FILE_OPEN_ERROR, http.StatusInternalServerError)

		}
	}
	defer feedsFile.Close()
	feedsFileData, err := ioutil.ReadAll(feedsFile)
	if err != nil {
		fmt.Printf("Error while reading feeds file ERROR: %s\n", err.Error())
		return nil, models.NewDataServiceError(FILE_READ_ERROR, http.StatusInternalServerError)
	}
	return feedsFileData, nil
}

//openFileAndFillObject opens the file at filename and tries to stuff the contents into object to fill, returning an error if any occurs
func OpenFileAndFillObject(objectToFill interface{}, fileName string) error {
	feedsFileData, err := ReadBytesOrError(fileName)
	if err != nil {
		fmt.Printf("Error while reading feeds file ERROR: %s\n", err.Error())
		return models.NewDataServiceError(FILE_READ_ERROR, http.StatusInternalServerError)
	}
	err = json.Unmarshal(feedsFileData, objectToFill)
	if err != nil {
		fmt.Printf("Error While unmarshalling feeds object ERROR: %s\n", err.Error())
		return models.NewDataServiceError(FILE_READ_ERROR, http.StatusInternalServerError)
	}

	return nil
}

func FillObject(data []byte, objectToFill interface{}) error {
	err := json.Unmarshal(data, objectToFill)
	if err != nil {
		fmt.Printf("Error While unmarshalling feeds object ERROR: %s\n", err.Error())
		return models.NewDataServiceError(UNMARSHALL_ERROR, http.StatusInternalServerError)
	}
	return nil
}
