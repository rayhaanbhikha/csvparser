package csvparser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func stringPty(val string) *string {
	return &val
}

func Parse[T any](csvReader *csv.Reader, csvRow T) ([]T, error) {
	headerIndexMap := make(map[string]int)
	csvRowType := reflect.TypeOf(csvRow)
	csvRows := make([]T, 0)

	if csvRowType == nil {
		return nil, errors.New("csvRow cannot be nil")
	}

	if csvRowType.Kind() != reflect.Struct {
		return nil, errors.New("csvRow must be a struct")
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(headerIndexMap) == 0 {
			for i, headerCol := range row {
				// what happens if we have 2 indexes pointing to the same column?
				headerIndexMap[headerCol] = i
			}
			continue
		}

		csvRowPtr := reflect.New(csvRowType)

		for i := 0; i < csvRowType.NumField(); i++ {
			structField := csvRowType.Field(i)
			if !structField.IsExported() {
				continue
			}

			omitEmptyFlagSet := false
			columnName := structField.Name
			if structField.Tag != "" {
				columnName, omitEmptyFlagSet = parseTag(structField.Tag)
			}

			colIndex, ok := headerIndexMap[columnName]
			if !ok || (colIndex > len(row)) {
				continue
			}

			colVal := row[colIndex]
			if colVal == "" && omitEmptyFlagSet {
				continue
			}

			csvRowPtrField := csvRowPtr.Elem().Field(i)
			if !csvRowPtrField.CanSet() {
				continue
			}

			// TODO: should we have additional safe guards before attempting to set the value?
			if csvRowPtrField.Kind() == reflect.Pointer {
				if colVal == "" {
					continue
				}
				clv := reflect.ValueOf(stringPty(colVal))
				csvRowPtrField.Set(clv)
			} else {
				csvRowPtrField.SetString(colVal)
			}
		}
		v, ok := csvRowPtr.Elem().Interface().(T)
		if !ok {
			return nil, fmt.Errorf("failed to map csvRow to type %T", csvRowType)
		}
		csvRows = append(csvRows, v)
	}

	return csvRows, nil
}

func parseTag(tag reflect.StructTag) (string, bool) {
	extractedTag := tag.Get("csv_header")
	name, omitEmpty, ok := strings.Cut(extractedTag, ",")
	if !ok || omitEmpty != "omitempty" {
		return name, false
	}
	return name, true
}
