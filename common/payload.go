package com

import (
	"bytes"
	"io"
)


type Payload struct {
	payloadName string //uniq identifier of the payload
	parted [][]byte
	resLen int
	res []byte

	fullLength int

}

// This Method Will Write Payload With Host
// Can be used In client side
func (pl *Payload)  WriteTo(writer io.Writer, host string) error {
	var err error
	for i, pload := range pl.parted {
		_, err = writer.Write(pload)
		if err != nil {
			return err
		}
		if i+1 == len(pl.parted) {
			break
		}
		_, err = writer.Write([]byte(host))
		if err !=nil {
			return err
		}
	}
	
	return nil
}

func (pl *Payload)  ReadRes(conn io.Reader) (int, error) {
	res := make([]byte, pl.resLen) //FIXME: can be use sync.Pool
	return io.ReadFull(conn, res)
}


func (pl *Payload) Name() string {
	return pl.payloadName
}

// This Method Will Write Response According to it's Payload
// Can be used In server side
// response won't be larger than 100KB so it's all right to wright whole response at once 
func (pl *Payload)  WriteRes(writer io.Writer) (int, error) {
	return writer.Write(pl.res)
	
}

func (pl *Payload) ReadAfterFirst(reader io.Reader) (string, error) {
	if len(pl.parted) == 1 {
		return "", nil
	}
	todetect := pl.parted[1]
	if len(todetect) > 5 {
		todetect = todetect[:5]
	}
	totalRead := 0
	host := []byte{}
	ch := make([]byte, 1)
	hostLen := 0
	
	// TODO: may be use bufio buffred reader
	for {
		n, err := reader.Read(ch)
		if err != nil {
			return "", err
		}
		totalRead += n
		
		if ch[0] == todetect[0] {
			to := make([]byte, len(todetect)-1)
			io.ReadFull(reader, to)
			if bytes.Equal(to, todetect[1:]) {
				hostLen = len(host)
				break
			} else {
				reader = NewBufReader(to, reader)
			}

		}
		host = append(host, ch[0])
	}

	reader = UnwrapReader(reader)

	shouldReadHostCount := len(pl.parted)-2 //already read one
	if shouldReadHostCount < 0 {
		shouldReadHostCount = 0
	}

	toread := (pl.fullLength - (len(pl.parted[0]) + len(todetect))) + (shouldReadHostCount * hostLen)
	ch = make([]byte, toread)
	for {
		n, err := reader.Read(ch)
		ch = ch[n:]
		if err != nil {
			return string(host), err
		}
		if len(ch) == 0 {
			return string(host), nil
		}
	}
}





























type PayloadOnce struct {
	payloadName string //uniq identifier of the payload

	parted [][]byte
	host []byte
	hostoff int
	copyhost bool
}

func (pl *PayloadOnce)  Read(p []byte) (int, error) {
	write := 0
	
	if pl.copyhost {
		now := copy(p, pl.host[pl.hostoff:])
		pl.hostoff += now

		if len(pl.host[pl.hostoff:]) > 0 {
			return len(p), nil
		}

		write += now
		pl.hostoff = 0
		pl.copyhost = false
		if now == len(p) {
			return len(p), nil
		}
	}

	for len(pl.parted) > 0 {
		now := copy(p[write:], pl.parted[0])
		write += now
		pl.parted[0] = pl.parted[0][now:]
		if len(pl.parted[0]) > 0 { //dst over
			return write, nil
		}
		if len(pl.parted) == 1 {
			break
		}

		pl.copyhost = true
		pl.parted = pl.parted[1:]

		pl.hostoff = copy(p[write:], pl.host)
		write += pl.hostoff
		if pl.hostoff != len(pl.host) {
			return write, nil
		}
		pl.copyhost = false
		pl.hostoff = 0
	}
	
	return 0, io.EOF
}


func (pl *PayloadOnce) Name() string {
	return pl.payloadName
}
