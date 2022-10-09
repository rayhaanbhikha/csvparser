package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func Parse[T any](csvReader *csv.Reader, csvRowMapping T) ([]T, error) {
	rowMappingVal := reflect.ValueOf(csvRowMapping)

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

		v, ok := csvRowPtr.Elem().Interface().(T)
		if !ok {
			return nil, fmt.Errorf("failed to map csvRow to type %T", rowMappingVal.Type())
		}

		csvRows = append(csvRows, v)
	}

	return csvRows, nil
}
