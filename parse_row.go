package csvparser

import (
	"reflect"
	"strings"
)

type csvRowMapper struct {
	fields  map[int]int
	rowType reflect.Type
}

func newCSVRowType(rowType reflect.Type) *csvRowMapper {
	return &csvRowMapper{
		rowType: rowType,
		fields:  make(map[int]int),
	}
}

func (c *csvRowMapper) addField(index, csvRowIndex int) {
	c.fields[index] = csvRowIndex
}

func (c *csvRowMapper) generate(row []string) reflect.Value {
	csvRowPtr := reflect.New(c.rowType)

	for i, csvRowIndex := range c.fields {
		csvRowPtrField := csvRowPtr.Elem().Field(i)
		colVal := row[csvRowIndex]
		setVal(&csvRowPtrField, colVal)
	}

	return csvRowPtr
}

func setVal[T any](rv *reflect.Value, actualVal T) {
	val := reflect.ValueOf(actualVal)
	if !rv.CanSet() {
		return
	}

	if rv.Kind() == reflect.Pointer {
		if val.IsZero() {
			return
		}
		rv.Set(reflect.ValueOf(&actualVal))
		return
	}
	rv.Set(val)
}

func rowMapper(rowVal reflect.Value, headers *csvHeaders) (*csvRowMapper, error) {
	if !rowVal.IsValid() {
		return nil, ErrCSVRowHasInvalidValue
	}

	underlyingStruct, err := extractStruct(rowVal)
	if err != nil {
		return nil, ErrCSVRowMustBeAStruct
	}
	underlyingStructType := underlyingStruct.Type()

	csvRowType := newCSVRowType(underlyingStructType)

	for i := 0; i < underlyingStructType.NumField(); i++ {
		structField := underlyingStructType.Field(i)

		if !structField.IsExported() {
			continue
		}

		if !isKind(structField.Type, reflect.String) {
			return nil, ErrCSVRowMappingInvalidField
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

func parseTag(tag reflect.StructTag) string {
	extractedTag := tag.Get("csv_header")
	name, _, _ := strings.Cut(extractedTag, ",")
	return name
}
