package csvparser

import (
	"errors"
	"reflect"
	"strings"
)

type csvRowType struct {
	fields  map[int]int
	rowType *reflect.Type
}

func newCSVRowType() *csvRowType {
	return &csvRowType{
		fields: make(map[int]int),
	}
}

func (c *csvRowType) addField(index, csvRowIndex int) {
	c.fields[index] = csvRowIndex
}

func (c *csvRowType) generate(rowType reflect.Type, row []string) reflect.Value {
	csvRowPtr := reflect.New(rowType)

	for i, csvRowIndex := range c.fields {
		csvRowPtrField := csvRowPtr.Elem().Field(i)
		if !csvRowPtrField.CanSet() {
			continue
		}

		colVal := row[csvRowIndex]
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

	return csvRowPtr
}

func isKind(t reflect.Type, k2 reflect.Kind) bool {
	switch k1 := t.Kind(); k1 {
	case k2:
		return true
	case reflect.Pointer:
		return isKind(t.Elem(), k2)
	default:
		return false
	}
}

func parseType(rowVal reflect.Value, headers *csvHeaders) (*csvRowType, error) {
	if !rowVal.IsValid() {
		return nil, ErrCSVRowHasInvalidValue
	}

	rowType := rowVal.Type()

	csvRowType := newCSVRowType()
	// should we allow pointers to a struct?
	if rowType.Kind() != reflect.Struct {
		return nil, ErrCSVRowMustBeAStruct
	}

	for i := 0; i < rowType.NumField(); i++ {
		structField := rowType.Field(i)

		if !structField.IsExported() {
			continue
		}

		if !isKind(structField.Type, reflect.String) {
			return nil, errors.New("expected field type to be string or *string")
		}

		columnName := structField.Name
		if structField.Tag != "" {
			columnName = parseTag(structField.Tag)
		}

		colIndex, ok := headers.get(columnName)
		if !ok {
			continue
		}

		csvRowType.addField(i, colIndex)
	}

	return csvRowType, nil
}

//func parseRow(csvRowType reflect.Type, csvHeaders *csvHeaders, row []string) (reflect.Value, error) {
//	csvRowPtr := reflect.New(csvRowType)
//
//	for i := 0; i < csvRowType.NumField(); i++ {
//		structField := csvRowType.Field(i)
//		if !structField.IsExported() {
//			continue
//		}
//
//		omitEmptyFlagSet := false
//		columnName := structField.Name
//		if structField.Tag != "" {
//			columnName, omitEmptyFlagSet = parseTag(structField.Tag)
//		}
//
//		colIndex, ok := csvHeaders.get(columnName)
//		if !ok || (colIndex > len(row)) {
//			continue
//		}
//
//		colVal := row[colIndex]
//		if colVal == "" && omitEmptyFlagSet {
//			continue
//		}
//
//		csvRowPtrField := csvRowPtr.Elem().Field(i)
//		if !csvRowPtrField.CanSet() {
//			continue
//		}
//
//		valToSet := reflect.ValueOf(colVal)
//		// TODO: should we have additional safe guards before attempting to set the value?
//		if csvRowPtrField.Kind() == reflect.Pointer {
//			if colVal == "" {
//				continue
//			}
//			valToSet = reflect.ValueOf(&colVal)
//		}
//		csvRowPtrField.Set(valToSet)
//	}
//	return csvRowPtr, nil
//}

func parseTag(tag reflect.StructTag) string {
	extractedTag := tag.Get("csv_header")
	name, _, _ := strings.Cut(extractedTag, ",")
	return name
}
