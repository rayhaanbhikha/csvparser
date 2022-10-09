package csvparser

import "encoding/csv"

type csvHeaders struct {
	headersByIndex map[string]int
}

func newCSVHeaders(csvReader *csv.Reader) (*csvHeaders, error) {
	headersByIndex := make(map[string]int)
	headers, err := csvReader.Read()
	if err != nil {
		return nil, err
	}
	for i, headerCol := range headers {
		// what happens if we have 2 indexes pointing to the same column?
		headersByIndex[headerCol] = i
	}
	return &csvHeaders{headersByIndex: headersByIndex}, err
}

func (c *csvHeaders) get(header string) (int, bool) {
	index, ok := c.headersByIndex[header]
	return index, ok
}
