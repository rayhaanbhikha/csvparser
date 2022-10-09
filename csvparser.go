package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func Parse[T any](csvReader *csv.Reader, csvRowMapping T) ([]T, error) {
	rowMappingVal := reflect.ValueOf(csvRowMapping)
	isPointer := rowMappingVal.Kind() == reflect.Pointer
	csvRows := make([]T, 0)
	csvHeaders, err := newCSVHeaders(csvReader)
	if err == io.EOF {
		return csvRows, nil
	}
	if err != nil {
		return nil, err
	}

	csvRowMapper, err := rowMapper(rowMappingVal, csvHeaders)
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

		csvRowPtr := csvRowMapper.generate(row)

		if !isPointer {
			csvRowPtr = csvRowPtr.Elem()
		}

		v, ok := csvRowPtr.Interface().(T)
		if !ok {
			return nil, fmt.Errorf("failed to map from %T to %T", csvRowPtr.Elem().Interface(), rowMappingVal.Type())
		}

		csvRows = append(csvRows, v)
	}

	return csvRows, nil
}
