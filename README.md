# csvparser

This library aims to parse CSV files and map a row to the provided struct.

```Go
package main

import (
    "encoding/csv"
    "fmt"
    "github.com/rayhaanbhikha/csvparser"
    "log"
    "strings"
)

var csvData = strings.NewReader(`name,age,gender
john,30,male
Rob,40,male
victoria,25,female
lizzy,,
alicia,,female`)

type Person struct {
	Name string `csv_header:"name"`
	Age *string `csv_header:"age"`
	Gender string `csv_header:"gender"`
}

func main() {
    csvReader := csv.NewReader(csvData)
	
    people, err := csvparser.Parse(csvReader, &Person{})
    if err != nil {
        log.Fatal(err)
    }
	
    fmt.Println(people)
	
    /* should expect the following
        []*Person{
            {Name: "john", Age: pointy("30"), Gender: "male"},
            {Name: "Rob", Age: pointy("40"), Gender: "male"},
            {Name: "victoria", Age: pointy("25"), Gender: "female"},
            {Name: "lizzy", Gender: ""},
            {Name: "alicia", Gender: "female"},
        }
    */
}
```

Additionally, you can also use the `ParseChan` method.

```Go
package main

import (
	"encoding/csv"
	"fmt"
	"github.com/rayhaanbhikha/csvparser"
	"log"
	"strings"
)

var csvData = strings.NewReader(`name,age,gender
john,30,male
Rob,40,male
victoria,25,female
lizzy,,
alicia,,female`)

type Person struct {
    Name string `csv_header:"name"`
    Age *string `csv_header:"age"`
    Gender string `csv_header:"gender"`
}

func main() {
    csvReader := csv.NewReader(csvData)
	
    parsedResultChan := csvparser.ParseChan[*Person](csvReader, &Person{})
    people := make([]*Person, 0)
	
    for parsedResult := range parsedResultChan {
        if parsedResult.Error != nil {
            log.Fatal(parsedResult.Error)
        }
        people = append(people, parsedResult.Row)
    }
	
	fmt.Println(people)
    /* should expect the following
        []*Person{
            {Name: "john", Age: pointy("30"), Gender: "male"},
            {Name: "Rob", Age: pointy("40"), Gender: "male"},
            {Name: "victoria", Age: pointy("25"), Gender: "female"},
            {Name: "lizzy", Gender: ""},
            {Name: "alicia", Gender: "female"},
        }
    */
}
```