package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	com "github.com/sadeepa24/netshoot/common"
)

type PayloadINfoFile struct {
	Format   string       `json:"format"`
	Output   string       `json:"output"`
	Payloads []OnePayload `json:"payloads"`
}

type OnePayload struct {
	Name     string `json:"name"`
	Skip     bool   `json:"skip"`
	Payload  string `json:"payload"`
	Response string `json:"response"`
}


func ReadPayloadInfo(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var PfileInfo PayloadINfoFile

	err = json.Unmarshal(file, &PfileInfo)
	if err != nil {
		fmt.Println("marshel err")
		return err
	}

	var (
		allpayloads, res [][]byte
		names []string
	)

	for _, payload := range PfileInfo.Payloads {
		if payload.Skip {
			continue
		}
		var dstP, dstR []byte
		switch PfileInfo.Format {
		case "hex":
			dstP, err = hex.DecodeString(payload.Payload)
			if err != nil {
				fmt.Println("err from here")
				return err
			}
			dstR, err = hex.DecodeString(payload.Response)
			if err != nil {
				return err
			}
		case "raw":
			dstP = []byte(payload.Payload)
			dstR = []byte(payload.Response)
		case "base64":
			return errors.New("not aqvailable yet")
		}
		allpayloads = append(allpayloads, dstP)
		res = append(res, dstR)
		names = append(names, payload.Name)
	}

	endFile, err := os.OpenFile(PfileInfo.Output, os.O_CREATE|os.O_RDWR, 0644)
	
	if err != nil {return err}
	endFile.Truncate(0)
	return com.CreatePayloadFile(allpayloads, res, names, endFile)
}