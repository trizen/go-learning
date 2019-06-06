package main

import (
	"fmt"
	"time"
)

// Array -> Hash -> Array -> Hash
// [{key => [{key => string}]}, {key => [{key => string}]}]

func main() {

	type Item struct {
		name, exec, icon, terminal string
	}
	var schema_data []map[string][]Item

	for i := 1; i < 12; i++ {

		schema_2 := make(map[string][]Item)
		m := fmt.Sprintf("%s", time.Month(i))

		schema_2[m] = append(schema_2[m], Item{
			"A name",
			"Exec",
			"Icon",
			"TerminaL",
		})
		schema_2[m] = append(schema_2[m], Item{
			fmt.Sprintf("Name: %d", i),
			"Exec=2",
			"Icon-2",
			"Terminal2",
		})
		schema_data = append(schema_data, schema_2)
	}

	for _, key := range schema_data {
		for _, value := range key {
			for _, data := range value {
				fmt.Println(data.name)
				fmt.Printf("%+v\n", data)
			}
		}
	}

	fmt.Printf("%+v\n", schema_data)
}
