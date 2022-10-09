package csvparser

//func ParseSlice[T any](csvReader *csv.Reader, csvRow *[]T) error {
//	var csvHeaders *csvHeaders
//
//	csvRows := make([]T, 0)
//
//	for {
//		row, err := csvReader.Read()
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			return err
//		}
//
//		if csvHeaders == nil {
//			csvHeaders = newCSVHeaders(row)
//			continue
//		}
//
//		csvRowPtr, err := parseRow(csvRowType, csvHeaders, row)
//		if err != nil {
//			return nil, err
//		}
//
//		v, ok := csvRowPtr.Elem().Interface().(T)
//		if !ok {
//			return nil, fmt.Errorf("failed to map csvRow to type %T", csvRowType)
//		}
//
//		csvRows = append(csvRows, v)
//	}
//
//	return csvRows, nil
//}
