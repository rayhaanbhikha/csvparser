package csvparser_test

import (
	"encoding/csv"
	"fmt"
	"github.com/rayhaanbhikha/csvparser"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

// FIXME: use strings.NewReader instead of opening files.
func TestParse(t1 *testing.T) {

	t1.Run("should return expected rows", func(t *testing.T) {
		type Row struct {
			Name   string  `csv_header:"name"`
			Age    *string `csv_header:"age"`
			Gender string  `csv_header:"gender"`
		}

		csvReader := csv.NewReader(mockCSVData())
		got, err := csvparser.Parse[Row](csvReader, Row{})
		require.NoError(t, err)
		expected := []Row{
			{Name: "john", Age: pointy("30"), Gender: "male"},
			{Name: "Rob", Age: pointy("40"), Gender: "male"},
			{Name: "victoria", Age: pointy("25"), Gender: "female"},
			{Name: "lizzy"},
			{Name: "alicia", Gender: "female"},
		}

		assert.ElementsMatch(t, expected, got)
	})

	t1.Run("should return expected rows", func(t *testing.T) {

		type Row struct {
			Name   string  `csv_header:"name"`
			Age    *string `csv_header:"age"`
			Gender string  `csv_header:"gender"`
		}

		csvReader := csv.NewReader(mockCSVData())
		got, err := csvparser.Parse[*Row](csvReader, &Row{})
		require.NoError(t, err)

		expected := []*Row{
			{Name: "john", Age: pointy("30"), Gender: "male"},
			{Name: "Rob", Age: pointy("40"), Gender: "male"},
			{Name: "victoria", Age: pointy("25"), Gender: "female"},
			{Name: "lizzy"},
			{Name: "alicia", Gender: "female"},
		}

		assert.ElementsMatch(t, expected, got)

	})

	t1.Run("should return error if v is an empty interface", func(t *testing.T) {
		var v any
		csvReader := csv.NewReader(mockCSVData())
		_, err := csvparser.Parse(csvReader, &v)
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is an empty interface", func(t *testing.T) {
		var v any
		csvReader := csv.NewReader(mockCSVData())
		_, err := csvparser.Parse(csvReader, &v)
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is a string", func(t *testing.T) {
		csvReader := csv.NewReader(mockCSVData())
		res, err := csvparser.Parse(csvReader, "some string")
		fmt.Println(res)
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is an invalid value", func(t *testing.T) {
		var c any

		csvReader := csv.NewReader(mockCSVData())
		_, err := csvparser.Parse(csvReader, c)
		require.Error(t, err, csvparser.ErrCSVRowHasInvalidValue)
	})

	t1.Run("should return error if v is set to nil", func(t *testing.T) {
		csvReader := csv.NewReader(mockCSVData())
		_, err := csvparser.Parse[any](csvReader, nil)
		require.Error(t, err, csvparser.ErrCSVRowHasInvalidValue)
	})
}
