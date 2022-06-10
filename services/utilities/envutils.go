package utilities

import (
	"fmt"
	"os"
)

func GetEnvAndLogOrError(envName string) (string, error) {
	envVal, err := GetEnvOrError(envName, fmt.Sprintf("%s is required", envName))
	if err != nil {
		return envVal, err
	}
	fmt.Printf("using %s: %s\n", envName, envVal)
	return envVal, nil
}

func GetEnvDefaultAndLog(envName, defaultVal string) string {
	envVal := GetEnvStringOrDefault(envName, defaultVal)
	fmt.Printf("using %s: %s\n", envName, envVal)
	return envVal
}

func GetEnvStringOrDefault(envName, defaultValue string) string {
	if val, exists := os.LookupEnv(envName); exists && len(val) > 0 {
		return val
	}
	return defaultValue
}

func GetEnvOrError(envName, errorMsg string) (string, error) {
	if val, exists := os.LookupEnv(envName); exists && len(val) > 0 {
		return val, nil
	}
	return "", fmt.Errorf(errorMsg)
}
