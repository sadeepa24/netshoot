package server

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
)


func TestServer(t *testing.T) {
	lgr, _ := zap.NewDevelopment()
	server, err := NewPayloadServer(context.Background(), lgr, config.PayloadServer{
		PayloadFile: "payload.dt",
		Ls: config.LsConfig{
			ListenAddr: "127.0.0.1:24000",
			Tls: config.TlsServer{
				Enabled: false,
			},
		},

	})
	if err != nil {
		t.Error(err)
	}


	//making test payloads and replace with server's payload object
	testbuf := &bytes.Buffer{}
	err = com.CreatePayloadFile(
		[][]byte{
			[]byte("h" +com.PayloadDelim+" ellllllssssssssssssssssssssssssssssssssssl"),
			[]byte("nope"),
			[]byte("Second "+com.PayloadDelim+"Payload"),
		},
		[][]byte{
			[]byte("First Payload Response"),
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
	allp, err := com.ReadPayload(testbuf)
	if err != nil {
		t.Error(err)
	}
	server.payloads  = allp
	server.payloadFirstSorted = allp.FirstPart()
	
	go server.Start()


	// real test
	tbuf := &bytes.Buffer{}
	tbuf.Write([]byte("h" +"--first test host--"+" ellllllssssssssssssssssssssssssssssssssssl"))

	
	server.handleconn(&serverTestConn{
		buf: tbuf,
	})

	tbuf.Reset()
	tbuf.Write([]byte("Second "+"-- second test host ---"+"Payload"))
	

	server.handleconn(&serverTestConn{
		buf: tbuf,
	})

	tbuf.Reset()
	tbuf.Write([]byte("nope"))
	

	server.handleconn(&serverTestConn{
		buf: tbuf,
	})

}

func TestPayloadFIrstDetect(t *testing.T) {
	server := PayloadServer{
		payloadFirstSorted: [][]byte{
			[]byte("hello"),
			[]byte("second Payload Longer"),
			[]byte("third payload much longer"),
		},
	}
	server.bytepool = sync.Pool{
		New: func() any {
			return make([]byte, len(server.payloadFirstSorted[len(server.payloadFirstSorted)-1]))
		},
	}

	testconn := &serverTestConn{
		buf: &bytes.Buffer{},
	}
	testconn.buf.Write([]byte("third payload much longer"))
	payloadnum, err := server.detectPayloadStrict(testconn)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(payloadnum)
}





type serverTestConn struct {
	buf *bytes.Buffer
}

func (s *serverTestConn) Read(b []byte) (n int, err error) {
	return s.buf.Read(b)
}

func (s *serverTestConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (s *serverTestConn) Close() error {
	return nil
}
























func (s *serverTestConn) LocalAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
}

func (s *serverTestConn) RemoteAddr() net.Addr {
	return &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)}
}

func (s *serverTestConn) SetDeadline(t time.Time) error {
	return nil
}

func (s *serverTestConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *serverTestConn) SetWriteDeadline(t time.Time) error {
	return nil
}

