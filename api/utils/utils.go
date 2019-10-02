package utils

// GetInteger - checks int pointer
func GetInteger(intPointer *int, defaultVal int) *int {
	if intPointer == nil {
		return &defaultVal
	}
	return intPointer
}
