package csvparser_test

import (
	"encoding/csv"
	"fmt"
	"github.com/rayhaanbhikha/csvparser"
	"testing"
)

func BenchmarkParse(b *testing.B) {
	csvReader := csv.NewReader(mockCSVData())

	type Row struct {
		Name   string  `csv_header:"name"`
		Age    *string `csv_header:"age"`
		Gender string  `csv_header:"gender"`
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := csvparser.Parse(csvReader, Row{})
		if err != nil {
			fmt.Println(err)
			b.FailNow()
		}
	}
	b.StopTimer()
}
