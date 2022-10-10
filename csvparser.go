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

func ParseChan[T any](csvReader *csv.Reader, v T) <-chan ParseChanResult[T] {
	resultChan := make(chan ParseChanResult[T])

	go func() {
		defer close(resultChan)
		rv := reflect.ValueOf(v)
		isPointer := rv.Kind() == reflect.Pointer

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

		csvRowMapper, err := rowMapper(rv, csvHeaders)
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

			mappedCSVRowPtr := csvRowMapper.generate(row)

			if !isPointer {
				mappedCSVRowPtr = mappedCSVRowPtr.Elem()
			}

			v, ok := mappedCSVRowPtr.Interface().(T)
			if !ok {
				resultChan <- ParseChanResult[T]{
					Error: fmt.Errorf("failed to map from %T to %T", mappedCSVRowPtr.Interface(), rv.Type()),
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
