package models

//DataServiceError defines a standard error type for any data service
type DataServiceError struct {
	Message string
	Code    int
}

//NewDataServiceError returns an error with the passed message and code
func NewDataServiceError(message string, code int) error {
	return DataServiceError{
		Message: message,
		Code:    code,
	}
}

func (err DataServiceError) Error() string {
	return err.Message
}

//ErrorCode returns the error code associated with this data service error
func (err DataServiceError) ErrorCode() int {
	return err.Code
}
