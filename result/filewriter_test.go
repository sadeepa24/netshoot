package result

import (
	"log"
	"os"
	"testing"
)

func TestFileWriter(t *testing.T) {
	file, err := os.OpenFile("testfile.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("open file failed")
	}

	writer := resultWriter{
		internal: file,
	}
	_, err = writer.Write([]byte("this is teststtsts  hshgdhsgd"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	_, err = writer.Write([]byte("<<<<<<<NEXTNEXT>>>>>>>"))
	writer.Close()
}