package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
)

func Parse[T any](csvReader *csv.Reader, csvRowMapping T) ([]T, error) {
	csvRows := make([]T, 0)
	for parsedResult := range ParseChan(csvReader, csvRowMapping) {
		if parsedResult.Error != nil {
			return nil, parsedResult.Error
		}
		csvRows = append(csvRows, parsedResult.Row)
	}
	return csvRows, nil
}

type ParseChanResult[T any] struct {
	Row   T
	Error error
}

func ParseChan[T any](csvReader *csv.Reader, csvRowMapping T) <-chan ParseChanResult[T] {
	resultChan := make(chan ParseChanResult[T])

	go func() {
		defer close(resultChan)
		rowMappingVal := reflect.ValueOf(csvRowMapping)
		isPointer := rowMappingVal.Kind() == reflect.Pointer

		csvHeaders, err := newCSVHeaders(csvReader)
		if err == io.EOF {
			return
		}
		if err != nil {
			resultChan <- ParseChanResult[T]{
				Error: err,
			}
			return
		}

		csvRowMapper, err := rowMapper(rowMappingVal, csvHeaders)
		if err != nil {
			resultChan <- ParseChanResult[T]{
				Error: err,
			}
			return
		}

		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				resultChan <- ParseChanResult[T]{
					Error: err,
				}
				continue
			}

			csvRowPtr := csvRowMapper.generate(row)

			if !isPointer {
				csvRowPtr = csvRowPtr.Elem()
			}

			v, ok := csvRowPtr.Interface().(T)
			if !ok {
				resultChan <- ParseChanResult[T]{
					Error: fmt.Errorf("failed to map from %T to %T", csvRowPtr.Elem().Interface(), rowMappingVal.Type()),
				}
				continue
			}

			resultChan <- ParseChanResult[T]{
				Row: v,
			}
		}

		return
	}()

	return resultChan
}
