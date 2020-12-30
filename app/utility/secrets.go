package utility

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func readSecretValueFromFile(filePath string) (secret string, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}
	if !fileInfo.Mode().IsRegular() {
		return "", fmt.Errorf("path to secret is not a file: %s", filePath)
	}
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return strings.TrimSpace(string(buffer)), nil
}
