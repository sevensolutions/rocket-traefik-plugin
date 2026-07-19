package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Expands the environment variable if it is enclosed in ${}. If the variable is not present, the original value is returned.
// Also supports reading the value from a file via ${file:/path/to/file}, eg. to consume a Docker/Kubernetes
// secret mounted as a file, since environment variables can be read by anyone with access to `docker inspect`.
func ExpandEnvironmentVariableString(value string) string {
	after, hasPrefix := strings.CutPrefix(value, "${")

	if hasPrefix {
		variableName, hasSuffix := strings.CutSuffix(after, "}")

		if hasSuffix {
			if filePath, isFile := strings.CutPrefix(variableName, "file:"); isFile {
				fileContent, err := os.ReadFile(filePath)

				if err == nil {
					return strings.TrimSpace(string(fileContent))
				}
			} else {
				variableValue, isDefined := os.LookupEnv(variableName)

				if isDefined {
					return variableValue
				}
			}
		}
	}

	return value
}

func ExpandEnvironmentVariableBoolean(value string, defaultValue bool) (bool, error) {
	after, hasPrefix := strings.CutPrefix(value, "${")

	if hasPrefix {
		variableName, hasSuffix := strings.CutSuffix(after, "}")

		if hasSuffix {
			variableValue, isDefined := os.LookupEnv(variableName)

			if isDefined {
				value = variableValue
			}
		}
	}

	if value == "true" || value == "1" {
		return true, nil
	} else if value == "false" || value == "0" {
		return false, nil
	} else if value != "" {
		return false, errors.New(fmt.Sprintf("Invalid boolean value \"%s\". Boolean values must be true/false or 1/0.", value))
	}

	return defaultValue, nil
}
