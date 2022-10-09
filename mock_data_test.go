package csvparser_test

import (
	"io"
	"strings"
)

func mockCSVData() io.Reader {
	return strings.NewReader(`name,age,gender
john,30,male
Rob,40,male
victoria,25,female
lizzy,,
alicia,,female`)
}
