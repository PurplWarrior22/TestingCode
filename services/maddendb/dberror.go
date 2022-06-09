package maddendb

//define a standard error type

type DbError struct {
	Message       string
	OriginalError error
}

//Error Interface Implementation
func (dbError *DbError) Error() string {
	return dbError.Message
}
