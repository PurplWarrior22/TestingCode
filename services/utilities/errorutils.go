package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PurplWarrior22/TestingCode/services/models"
)

//error handling utilities

//StatusCodeError checks for error type and returns an appropriate status code
func StatusCodeError(err error) int {
	switch errorType := err.(type) {
	case models.DataServiceError:
		return errorType.Code
	default:
		if code := seeIfCodeFieldExists(err); code != nil {
			return *code
		}
		return http.StatusInternalServerError
	}
}

//seeIfCodeFieldExists checks if we can extract any "code" field to
func seeIfCodeFieldExists(err error) *int {
	dataBytes, err := json.Marshal(err)
	if err != nil {
		return nil
	}
	codeField := struct {
		LowerCode *int `json:"code,omitempty"`
		UpperCode *int `json:"Code,omitempty"`
	}{}
	err = json.Unmarshal(dataBytes, &codeField)
	fmt.Println(codeField)
	if err != nil {
		return nil
	}
	if codeField.LowerCode != nil {
		return codeField.LowerCode
	}
	return codeField.UpperCode

}
