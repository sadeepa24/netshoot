package hostmanager

import (
	"fmt"
	"testing"

	config "github.com/sadeepa24/netshoot/configs"
)

//Don't think about thease test 🤣
func TestHostFile(t *testing.T) {
	file, err := newfile(config.Hostfile{
		Hostfile: "testfile",
		MaxConcurrent: 10,
	})
	if err != nil {
		t.Error(err)
	}

	file.initialize(55)
	
	if err != nil {
		t.Error(err)
	}

	for file.available() {
		printslice(file.next())
		fmt.Println()
	}


}

func printslice(s []string) {
	for _, i := range s {
		fmt.Println(i)
	}
}