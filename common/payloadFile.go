package com

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"sort"
)

// payload file structure

const MaxPayloadFileSize = 1024 * 1024 * 32
const PayloadDelim string = "<--netshoot-->" 

/*
paylaod file structure
	2				1				 X				4			4			X       X
(PayLoadCount)(payloadNameLength)(PayloadName)(PayloadLen)(ResponseLen)(Payload)(Response)

*/

// must return as sorted
func ReadPayloadFile(path string) (AllPayload, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	
	if err != nil {
		return nil, errors.New("payload file open error: " +err.Error())
	}
	defer file.Close()
	return ReadPayload(file)

}

func ReadPayload(reader io.Reader) (AllPayload, error)  {
	allPayloads, err := readAllPayload(reader)
	if err != nil {
		return nil, errors.New("payload file read err: " + err.Error())
	}
	sort.SliceStable(allPayloads, func(i, j int) bool {
		return len(allPayloads[i].parted[0]) < len(allPayloads[j].parted[0])
	})


	return AllPayload(allPayloads), nil
}

func CreatePayloadFile(payloads, response [][]byte, names []string, writer io.Writer) error {
	if len(payloads) != len(response) || len(payloads) != len(names) {
		return ErrNotEnogufpayloadOrRes
	}
	var err error
	if err = binary.Write(writer, binary.BigEndian, uint16(len(payloads))); err != nil {
		return err
	}

	for i := 0; i < len(payloads); i++ {
		if err = binary.Write(writer, binary.BigEndian, uint8(len(names[i]))); err != nil {
			return err
		}
		_, err = writer.Write([]byte(names[i]))
		if err != nil {
			return err
		}
		if err = binary.Write(writer, binary.BigEndian, uint32(len(payloads[i]))); err != nil {
			return err
		}
		if err = binary.Write(writer, binary.BigEndian, uint32(len(response[i]))); err != nil {
			return err
		}
		_, err = writer.Write([]byte(payloads[i]))
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(response[i]))
		if err != nil {
			return err
		}
	}
	return nil
}

var ErrNotEnogufpayloadOrRes = errors.New("payload and response count missmatch")


func readAllPayload(reader io.Reader) ([]Payload, error) {
	var payload []Payload
	
	var payloadCount uint16
	if err := binary.Read(reader, binary.BigEndian, &payloadCount); err != nil {
		return payload, err
	}

	for i := 0; i < int(payloadCount); i++ {
		pl, err := readOnePayload(reader) 
		if err != nil {
			return payload, err
		} 
		payload = append(payload, pl)

	}
	return payload, nil

} 

func readOnePayload(reader io.Reader) (Payload, error) {
	var payload Payload

	// Read payloadNameLength (1 byte)
	var payloadNameLength byte
	if err := binary.Read(reader, binary.BigEndian, &payloadNameLength); err != nil {
		return payload, err
	}

	// Read PayloadName (X bytes)
	payloadName := make([]byte, payloadNameLength)
	if _, err := io.ReadFull(reader, payloadName); err != nil {
		return payload, err
	}
	payload.payloadName = string(payloadName)

	// Read PayloadLen (4 bytes)
	var payloadLen uint32
	if err := binary.Read(reader, binary.BigEndian, &payloadLen); err != nil {
		return payload, err
	}

	// Read ResponseLen (4 bytes)
	var responseLen uint32
	if err := binary.Read(reader, binary.BigEndian, &responseLen); err != nil {
		return payload, err
	}

	payload.resLen = int(responseLen)


	// Read Payload (X bytes)
	payloadData := make([]byte, payloadLen)
	if _, err := io.ReadFull(reader, payloadData); err != nil {
		return payload, err
	}
	payload.parted = spiltPayload(payloadData)
	for _, r := range payload.parted {
		payload.fullLength += len(r)
	}
	
	responseData := make([]byte, responseLen)
	if _, err := io.ReadFull(reader, responseData); err != nil {
		return payload, err
	}

	payload.res = responseData




	return payload, nil
}

// this will run only starting it's allright to itrate all byte it's just one time
func spiltPayload(payload []byte) [][]byte {
	var end [][]byte
	var lastOccurance int
	if len(payload) >= len(PayloadDelim) {
		if bytes.Equal(payload[:len(PayloadDelim)],  []byte(PayloadDelim)) {
			payload = payload[len(PayloadDelim):]
		}
	}

	for i := 0; i < len(payload); i++ {
		if payload[i] == PayloadDelim[0] {

			if len(payload[i:]) >= len(PayloadDelim) {
				if bytes.Equal(payload[i:i+len(PayloadDelim)], []byte(PayloadDelim)) {
					if lastOccurance > 0 {
						lastOccurance += len(PayloadDelim)
					}
					if len(payload[lastOccurance:i]) > 0 {
						end = append(end, payload[lastOccurance:i])
					}
					lastOccurance = i
				}
			}
		}
	}
	if lastOccurance > 0 {
		end = append(end, payload[lastOccurance+len(PayloadDelim):])
	} else {
		end = append(end, payload)
	}
	
	return end
}


type AllPayload []Payload

//first part as sorted
func (a *AllPayload) FirstPart() [][]byte {
	ft := [][]byte{}
	for _, pload := range *a {
		ft = append(ft, pload.parted[0])
	}
	return ft
}