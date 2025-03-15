package main

import (
	"fmt"
	"testing"
)

func TestPgen(t *testing.T) {
	err := ReadPayloadInfo("test.json")
	fmt.Println(err)
}