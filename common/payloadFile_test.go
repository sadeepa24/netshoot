package com

import (
	"bytes"
	"fmt"
	"sort"
	"testing"
)

func TestPayloadSpiltRead(t *testing.T) {
	test := []byte(
		PayloadDelim + "s" + PayloadDelim + 
		"end of the test payload" + PayloadDelim + "wqdawedfawedfawe"+
		"s" + PayloadDelim +  PayloadDelim + 
		 PayloadDelim + "----------------- " + PayloadDelim,
	)

	sp := spiltPayload(test)
	
	for _, v := range sp {
		fmt.Println(string(v))
	}

}


func TestPayloadWrite(t *testing.T) {
	
	testbuf := &bytes.Buffer{}

	err := CreatePayloadFile(
		[][]byte{
			[]byte("this is the teststs paysd ewjfenfcljnecenfefweofow3efoweqjnf" +PayloadDelim+" ellllllssssssssssssssssssssssssssssssssssl"),
			[]byte("wefcawefvewrfvcews" +PayloadDelim+"drgvergvwergvwerdg" +PayloadDelim+"vwedrgvwer"),
			[]byte("Second Paygvwergvwerfgvwergvwergwer" +PayloadDelim+"gvwerdfytjyujkmuyk" +PayloadDelim+"myload"),
		},

		
		[][]byte{
			[]byte("First Payload"),
			[]byte("Nope Response"),
			[]byte("Second Response"),
		
		},

		[]string{
			"firstpayload",
			"small",
			"secondPayload",
		}, 


		testbuf,
	)
	if err != nil {
		t.Error(err)
	}

	// file, err := os.OpenFile("payload.dt", os.O_CREATE|os.O_WRONLY, 0644)
	// file.ReadFrom(testbuf)
	// return 

	payloads, err := readAllPayload(testbuf)
	if err != nil {
		t.Error(err)
	}
	sort.SliceStable(payloads, func(i, j int) bool {
		return len(payloads[i].parted[0]) < len(payloads[j].parted[0])
	})

	for _, v := range payloads {
		fmt.Println(v.payloadName)
		fmt.Println(v.resLen)
		fmt.Println(v.parted)
		fmt.Println(string(v.res))
		fmt.Println()
	}
}