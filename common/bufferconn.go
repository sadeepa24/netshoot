package com

import (
	"io"
	"net"
)

type BufConn struct {
	bufreader *BufReader
	net.Conn
}

func NewBufConn(buf []byte, conn net.Conn) *BufConn {
	return &BufConn{
		bufreader: NewBufReader(buf, conn),
		Conn: conn,
	}
}

func (bu *BufConn) Read(b []byte) (int, error) {
	return bu.bufreader.Read(b)
}


type BufReader struct {
	buf     []byte
	io.Reader
}
func NewBufReader(buf []byte, reader io.Reader) *BufReader {
	return &BufReader{
		buf: buf,
		Reader: reader,
	}
}

func (bu *BufReader) Read(b []byte) (int, error) {
	write := 0
	if len(bu.buf) > 0 {
		write = copy(b, bu.buf)
		bu.buf = bu.buf[write:]
	}
	n, err := bu.Reader.Read(b[write:])
	return write+n, err
}

func UnwrapConn(rd io.Reader) io.Reader {
	if nd, ok := rd.(*BufReader); ok {
		return UnwrapReader(nd.Reader)
	} else {
		return rd
	}
}

func UnwrapReader(rd io.Reader) io.Reader {
	if nd, ok := rd.(*BufReader); ok {
		return UnwrapReader(nd.Reader)
	} else {
		return rd
	}
}