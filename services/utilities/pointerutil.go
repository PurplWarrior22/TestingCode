package utilities

//StrPtr returns a pointer to a string whose contents are stringValue
func StrPtr(stringValue string) *string {
	return &stringValue
}

//IntPtr returns a pointer to an int with the content of intValue
func IntPtr(intValue int) *int {
	return &intValue
}
