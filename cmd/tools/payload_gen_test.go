package tools_test

import (
	"fmt"
	"testing"

	"github.com/sadeepa24/netshoot/cmd/tools"
)

func TestPgen(t *testing.T) {
	err := tools.GenPayloadFile("test.json")
	fmt.Println(err)
}