package csvparser_test

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

type csvRowStruct struct {
	Name   string  `csv_header:"name"`
	Age    *string `csv_header:"age"`
	Gender string  `csv_header:"gender"`
}

func Test_normalParse(t *testing.T) {
	csvReader := csv.NewReader(mockCSVData)
	results, err := normalParse(csvReader)
	require.NoError(t, err)

	expectedResults := []*csvRowStruct{
		{Name: "john", Age: stringPty("30"), Gender: "male"},
		{Name: "Rob", Age: stringPty("40"), Gender: "male"},
		{Name: "victoria", Age: stringPty("25"), Gender: "female"},
		{Name: "lizzy"},
		{Name: "alicia", Gender: "female"},
	}

	assert.ElementsMatch(t, expectedResults, results)
}

func normalParse(csvReader *csv.Reader) ([]*csvRowStruct, error) {
	headerMapIndex := make(map[string]int)
	csvRows := make([]*csvRowStruct, 0)
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(headerMapIndex) == 0 {
			for i, col := range row {
				headerMapIndex[col] = i
			}
			continue
		}

		csvRow := &csvRowStruct{}
		for colName, colIndex := range headerMapIndex {
			colVal := row[colIndex]
			switch colName {
			case "name":
				csvRow.Name = colVal
			case "age":
				if colVal != "" {
					csvRow.Age = &colVal
				}
			case "gender":
				csvRow.Gender = colVal
			}
		}
		csvRows = append(csvRows, csvRow)
	}

	return csvRows, nil
}
