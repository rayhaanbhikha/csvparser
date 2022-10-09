package csvparser_test

import "strings"

var mockCSVData = strings.NewReader(`name,age,gender
john,30,male
Rob,40,male
victoria,25,female
lizzy,,
alicia,,female`)
