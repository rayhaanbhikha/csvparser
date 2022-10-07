package csvparser_test

import (
	"encoding/csv"
	"github.com/rayhaanbhikha/csvparser"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func stringPty(val string) *string {
	return &val
}

type csvRowStruct struct {
	Name string  `csv_header:"name"`
	Age  *string `csv_header:"age"`
	Sex  string  `csv_header:"sex"`
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
				csvRow.Age = &colVal
			case "sex":
				csvRow.Sex = colVal
			}
		}
		csvRows = append(csvRows, csvRow)
	}

	return csvRows, nil
}

func BenchmarkParse(b *testing.B) {
	file, _ := os.Open("./data/sample_csv")

	csvReader := csv.NewReader(file)

	for i := 0; i < b.N; i++ {
		_, err := csvparser.Parse(csvReader, csvRowStruct{})
		if err != nil {
			b.FailNow()
		}
	}
}

func Benchmark_normalParse(b *testing.B) {
	file, _ := os.Open("./data/sample_csv")

	csvReader := csv.NewReader(file)

	for i := 0; i < b.N; i++ {
		_, err := normalParse(csvReader)
		if err != nil {
			b.FailNow()
		}
	}
}

func Test_normalParse(t *testing.T) {
	file, err := os.Open("./data/sample_csv")
	require.NoError(t, err)

	csvReader := csv.NewReader(file)
	results, err := normalParse(csvReader)
	require.NoError(t, err)

	expectedResults := []*csvRowStruct{
		{Name: "john", Age: stringPty("30"), Sex: "male"},
		{Name: "Rob", Age: stringPty("40"), Sex: "male"},
		{Name: "victoria", Age: stringPty("25"), Sex: "female"},
		{Name: "lizzy"},
		{Name: "alicia", Sex: "female"},
	}

	assert.ElementsMatch(t, expectedResults, results)
}

func TestParse(t *testing.T) {
	file, err := os.Open("./data/sample_csv")
	require.NoError(t, err)

	csvReader := csv.NewReader(file)
	results, err := csvparser.Parse[csvRowStruct](csvReader, csvRowStruct{})
	require.NoError(t, err)

	expectedResults := []csvRowStruct{
		{Name: "john", Age: stringPty("30"), Sex: "male"},
		{Name: "Rob", Age: stringPty("40"), Sex: "male"},
		{Name: "victoria", Age: stringPty("25"), Sex: "female"},
		{Name: "lizzy"},
		{Name: "alicia", Sex: "female"},
	}

	assert.ElementsMatch(t, expectedResults, results)
}
