package util

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"os"
)

func ReadTask(ctx context.Context, file string) (string, []string, []string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return "", nil, nil, err
	}
	reader := csv.NewReader(bytes.NewReader(content))
	records, err := reader.ReadAll()
	if err != nil {
	}

	if len(records) != 3 {
		return "", nil, nil, errors.New("format error 3 lines")
	}
	if len(records[0]) != 1 {
		return "", nil, nil, errors.New("format error line 1 need one elemenet")
	}
	if len(records[1]) < 1 {
		return "", nil, nil, errors.New("format error line 2 at least need one elemenet")
	}
	if len(records[2]) < 1 {
		return "", nil, nil, errors.New("format error line 3 at least need one elemenet")
	}
	sourceLanguage := records[0][0]
	destLanguage := records[1]
	contents := records[2]
	return sourceLanguage, destLanguage, contents, nil
}

func WriteResult(ctx context.Context, file string, results [][]string) error {

}
