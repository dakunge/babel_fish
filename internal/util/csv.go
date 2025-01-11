package util

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"os"

	"github.com/zeromicro/go-zero/core/logc"
)

func ReadTask(ctx context.Context, file string) (string, []string, []string, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		logc.Errorf(ctx, "read file %v, err %v", file, err)
		return "", nil, nil, err
	}
	reader := csv.NewReader(bytes.NewReader(content))
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		logc.Errorf(ctx, "read  contents %v, err %v", content, err)
		return "", nil, nil, err
	}

	logc.Debugf(ctx, "read file %v, records %v", file, records)
	if len(records) != 3 {
		logc.Errorf(ctx, "read reacords %v, format", records)
		return "", nil, nil, errors.New("format error 3 lines")
	}
	if len(records[0]) != 1 {
		logc.Errorf(ctx, "read reacords %v, format", records)
		return "", nil, nil, errors.New("format error line 1 need one elemenet")
	}
	if len(records[1]) < 1 {
		logc.Errorf(ctx, "read reacords %v, format", records)
		return "", nil, nil, errors.New("format error line 2 at least need one elemenet")
	}
	if len(records[2]) < 1 {
		logc.Errorf(ctx, "read reacords %v, format", records)
		return "", nil, nil, errors.New("format error line 3 at least need one elemenet")
	}
	sourceLanguage := records[0][0]
	destLanguage := records[1]
	contents := records[2]
	return sourceLanguage, destLanguage, contents, nil
}

func WriteResult(ctx context.Context, file string, results [][]string) error {
	//fmt.Println(records)
	f, err := os.Create(file) //创建文件
	if err != nil {
		logc.Errorf(ctx, "write results create file %v, err %v", file, err)
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	if err != nil {
		logc.Errorf(ctx, "write results write utf-8 bom  err %v", err)
		return err
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(results)
	if err != nil {
		logc.Errorf(ctx, "write results %v, err %v", results, err)
		return err
	}

	w.Flush()
	return nil
}
