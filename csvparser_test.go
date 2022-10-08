package csvparser_test

import (
	"encoding/csv"
	"github.com/rayhaanbhikha/csvparser"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type csvRowStruct struct {
	Name string  `csv_header:"name"`
	Age  *string `csv_header:"age"`
	Sex  string  `csv_header:"sex"`
}

func BenchmarkParse(b *testing.B) {
	file, _ := os.Open("./data/sample_csv")

	csvReader := csv.NewReader(file)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := csvparser.Parse(csvReader, csvRowStruct{})
		if err != nil {
			b.FailNow()
		}
	}
	b.StopTimer()
}

// FIXME: use strings.NewReader instead of opening files.
func TestParse(t1 *testing.T) {

	t1.Run("should return expected rows", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		type Row struct {
			Name string  `csv_header:"name"`
			Age  *string `csv_header:"age"`
			Sex  string  `csv_header:"sex"`
		}

		csvReader := csv.NewReader(file)
		got, err := csvparser.Parse[Row](csvReader, Row{})
		require.NoError(t, err)
		expected := []Row{
			{Name: "john", Age: pointy("30"), Sex: "male"},
			{Name: "Rob", Age: pointy("40"), Sex: "male"},
			{Name: "victoria", Age: pointy("25"), Sex: "female"},
			{Name: "lizzy"},
			{Name: "alicia", Sex: "female"},
		}
		assert.ElementsMatch(t, expected, got)
	})

	t1.Run("should return error if v is an empty interface", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		var v any

		csvReader := csv.NewReader(file)
		_, err = csvparser.Parse(csvReader, &v)
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is an empty interface", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		var v any

		csvReader := csv.NewReader(file)
		_, err = csvparser.Parse(csvReader, &v)
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is a string", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		csvReader := csv.NewReader(file)
		_, err = csvparser.Parse(csvReader, "some string")
		require.Error(t, err, csvparser.ErrCSVRowMustBeAStruct)
	})

	t1.Run("should return error if v is an invalid value", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		var c any

		csvReader := csv.NewReader(file)
		_, err = csvparser.Parse(csvReader, c)
		require.Error(t, err, csvparser.ErrCSVRowHasInvalidValue)
	})

	t1.Run("should return error if v is set to nil", func(t *testing.T) {
		file, err := os.Open("./data/sample_csv")
		require.NoError(t, err)
		t.Cleanup(func() {
			file.Close()
		})

		csvReader := csv.NewReader(file)
		_, err = csvparser.Parse[any](csvReader, nil)
		require.Error(t, err, csvparser.ErrCSVRowHasInvalidValue)
	})
}
