package main

import (
	"fmt"
	"strings"
)

var spaces int = 0     // the current number of spaces used for indentation
var spaces_num int = 4 // the number of spaces used to increment and decrement the `spaces` variable

type MyItem map[string]int
type MyStruct []MyItem
type MyKVStruct map[string]MyStruct

type JSON interface {
	encode_json() string
}

func to_json(t JSON) string {
	return t.encode_json()
}

func (p MyItem) encode_json() string {
	ind := strings.Repeat(" ", spaces)
	s := "{\n"

	spaces += spaces_num
	for k, v := range p {
		s += strings.Repeat(" ", spaces) + fmt.Sprintf("\"%s\": %d", k, v) + ",\n"
	}

	spaces -= spaces_num
	return s + ind + "}"
}

func (p MyStruct) encode_json() string {

	ind := strings.Repeat(" ", spaces)
	s := ind + "[\n"

	spaces += spaces_num
	for _, item := range p {
		if enc := to_json(item); enc != "" {
			s += strings.Repeat(" ", spaces) + enc + ",\n"
		}
	}

	spaces -= spaces_num
	return s + ind + "]"
}

func (p MyKVStruct) encode_json() string {

	ind := strings.Repeat(" ", spaces)
	s := "{\n"

	spaces += spaces_num
	for k, v := range p {
		if enc := fmt.Sprintf("\"%s\":\n%s", k, to_json(v)); enc != "" {
			s += strings.Repeat(" ", spaces) + enc + ",\n"
		}
	}

	spaces -= spaces_num
	return s + ind + "}"
}

func main() {

	// Create a simple struct
	s1 := make(MyStruct, 2)

	// Create and asssign some items
	s1[0] = make(MyItem)
	s1[0]["a"] = 1
	s1[1] = make(MyItem)
	s1[1]["b"] = 2

	fmt.Println("==> Array of maps:")
	fmt.Println(to_json(s1)) // call function to_json()

	// Create a nested struct
	s2 := make(MyKVStruct, 2)
	s2["my_key"] = s1

	s2["my_key"][0]["a"] = -42
	s2["my_key"][0]["a1"] = -42

	s2["my_key"][1]["b"] = -42
	s2["my_key"][1]["b1"] = -42

	s2["my_other_key"] = s1

	fmt.Println("==> Map of array of maps: ")
	fmt.Println(to_json(s2)) // call function to_json()
}
