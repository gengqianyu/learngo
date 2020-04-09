package main

import (
	"fmt"
	"testing"
)

func TestNilOrNot(t *testing.T) {
	var s *TestStruct
	if NilOrNot(s) == true {
		t.Errorf("expacted:%t,got:%t", false, true)
	}

}

func ExampleNilOrNot() {
	var s *TestStruct
	fmt.Println(s == nil)
	fmt.Println(NilOrNot(s))

	//output:
	//true
	//false

}
