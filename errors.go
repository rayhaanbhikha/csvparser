package csvparser

import "errors"

var (
	ErrCSVRowMustBeAStruct   = errors.New("csvRow value must be a struct")
	ErrCSVRowHasInvalidValue = errors.New("csvRow has invalid value")
)
