package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func Parse[T any](csvReader *csv.Reader, csvRow T) ([]T, error) {
	csvRowValue := reflect.ValueOf(csvRow)

	csvRows := make([]T, 0)
	csvHeaders, err := newCSVHeaders(csvReader)
	if err == io.EOF {
		return csvRows, nil
	}
	if err != nil {
		return nil, err
	}

	csvRowType, err := parseType(csvRowValue, csvHeaders)
	if err != nil {
		return nil, err
	}

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		csvRowPtr := csvRowType.generate(csvRowValue.Type(), row)

		v, ok := csvRowPtr.Elem().Interface().(T)
		if !ok {
			return nil, fmt.Errorf("failed to map csvRow to type %T", csvRowType)
		}

		csvRows = append(csvRows, v)
	}

	return csvRows, nil
}
