package com

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func TestPayload(t *testing.T) {
	testParted := [][]byte{
		[]byte("firsstssssspasyloasd"),
		[]byte("secondpayload"),
		[]byte("thirdpayload"),
	}
	testRes := []byte("test response")
	
	testPayload := Payload{
		payloadName: "testName",
		parted: testParted,
		resLen: len(testRes),
		res: testRes,
		fullLength: fulllength(testParted),
	}

	testwriter := &testRW{}
	
	testPayload.WriteTo(testwriter, "--s-s-s-ssssss-ss-ss-sssss-s")
	fmt.Println(string(testwriter.allrecived))

	mt := make([]byte, len(testParted[0]))

	testwriter.Read(mt)
	fmt.Println(string(testwriter.allrecived))

	host, err := testPayload.ReadAfterFirst(testwriter)
	fmt.Println("captured host: ", host)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(testwriter.allrecived))
	
}


func  fulllength(bt [][]byte) int {
	var m int
	for _, kk := range bt {
		m += len(kk)
	}
	return m
}


type testRW struct {
	allrecived []byte
}

func (t *testRW) Write(p []byte) (n int, err error) {
	t.allrecived = append(t.allrecived, p...)
	return len(p), nil
}

func (t *testRW) Read(p []byte) (int, error) {
	n := copy(p, t.allrecived)
	t.allrecived = t.allrecived[n:]
	if len(t.allrecived) == 0 {
		return n, io.EOF
	}
	return n, nil
}