package csvparser_test

import (
	"encoding/csv"
	"github.com/rayhaanbhikha/csvparser"
	"os"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	file, _ := os.Open("./data/sample_csv")

	csvReader := csv.NewReader(file)

	type Row struct {
		Name   string  `csv_header:"name"`
		Age    *string `csv_header:"age"`
		Gender string  `csv_header:"gender"`
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := csvparser.Parse(csvReader, Row{})
		if err != nil {
			b.FailNow()
		}
	}
	b.StopTimer()
}
