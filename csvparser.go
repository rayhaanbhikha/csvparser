package csvparser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type csvHeaders struct {
	headersByIndex map[string]int
}

var ErrCSVRowMustBeAStruct = errors.New("csvRow value must be a struct")
var ErrCSVRowHasInvalidValue = errors.New("csvRow has invalid value")

func newCSVHeaders(headers []string) *csvHeaders {
	headersByIndex := make(map[string]int)
	for i, headerCol := range headers {
		// what happens if we have 2 indexes pointing to the same column?
		headersByIndex[headerCol] = i
	}
	return &csvHeaders{headersByIndex: headersByIndex}
}

func (c *csvHeaders) get(header string) (int, bool) {
	index, ok := c.headersByIndex[header]
	return index, ok
}

func Parse[T any](csvReader *csv.Reader, csvRow T) ([]T, error) {
	csvRowValue := reflect.ValueOf(csvRow)

	// TODO: maybe we allow pointers?
	//if csvRowType.Kind() != reflect.Pointer {
	//	return nil, ErrCSVRowMustBeAPointer
	//}
	if !csvRowValue.IsValid() {
		return nil, ErrCSVRowHasInvalidValue
	}

	csvRowType := csvRowValue.Type()

	if csvRowType.Kind() != reflect.Struct {
		return nil, ErrCSVRowMustBeAStruct
	}

	var csvHeaders *csvHeaders
	csvRows := make([]T, 0)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if csvHeaders == nil {
			csvHeaders = newCSVHeaders(row)
			continue
		}

		csvRowPtr, err := parseRow(csvRowType, csvHeaders, row)
		if err != nil {
			return nil, err
		}

		v, ok := csvRowPtr.Elem().Interface().(T)
		if !ok {
			return nil, fmt.Errorf("failed to map csvRow to type %T", csvRowType)
		}

		csvRows = append(csvRows, v)
	}

	return csvRows, nil
}

func parseRow(csvRowType reflect.Type, csvHeaders *csvHeaders, row []string) (reflect.Value, error) {
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

		colIndex, ok := csvHeaders.get(columnName)
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

		valToSet := reflect.ValueOf(colVal)
		// TODO: should we have additional safe guards before attempting to set the value?
		if csvRowPtrField.Kind() == reflect.Pointer {
			if colVal == "" {
				continue
			}
			valToSet = reflect.ValueOf(&colVal)
		}
		csvRowPtrField.Set(valToSet)
	}
	return csvRowPtr, nil
}

func parseTag(tag reflect.StructTag) (string, bool) {
	extractedTag := tag.Get("csv_header")
	name, omitEmpty, ok := strings.Cut(extractedTag, ",")
	if !ok || omitEmpty != "omitempty" {
		return name, false
	}
	return name, true
}
