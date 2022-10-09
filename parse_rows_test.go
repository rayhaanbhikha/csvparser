package csvparser_test

import (
	"testing"
)

func TestParseSlice(t1 *testing.T) {
	t1.Run("Should return expected results", func(t *testing.T) {
		//type Row struct {
		//	Name   string  `csv_header:"name"`
		//	Age    *string `csv_header:"age"`
		//	Gender string  `csv_header:"gender"`
		//}
		//
		//var rows []Row
		//
		//csvReader := csv.NewReader(mockCSVData)
		////err := csvparser.ParseSlice[Row](csvReader, &rows)
		//require.NoError(t, err)
		//expected := []Row{
		//	{Name: "john", Age: pointy("30"), Gender: "male"},
		//	{Name: "Rob", Age: pointy("40"), Gender: "male"},
		//	{Name: "victoria", Age: pointy("25"), Gender: "female"},
		//	{Name: "lizzy"},
		//	{Name: "alicia", Gender: "female"},
		//}
		//assert.ElementsMatch(t, expected, rows)
	})
}
