package server

import (
	"crypto/tls"
	"net"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
)

type MixedListner struct {
	net.Listener
	tlsconf *tls.Config
}


func NewMixedLs(conf config.LsConfig) (*MixedListner, error) {
	ls := MixedListner{}

	if conf.Tls.Enabled {
		certs, err := tls.LoadX509KeyPair(conf.Tls.Cert, conf.Tls.Key)
		if err != nil {
			return nil, err
		}
		ls.tlsconf = &tls.Config{
			Certificates: []tls.Certificate{certs},
			InsecureSkipVerify: true,
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS13,
		}
	}
	
	laddr, err := net.ResolveTCPAddr("tcp", conf.ListenAddr)
	if err != nil {
		return nil, err
	}
	ls.Listener, err = net.ListenTCP("tcp", laddr)
	return &ls, err
}

func (m *MixedListner)  Accept() (net.Conn, error) {
	conn, err := m.Listener.Accept()
	if err != nil {
		return conn, err
	}

	if m.tlsconf != nil {		
		ft := make([]byte, 5)
		
		_, err := conn.Read(ft)
		if err != nil {
			return conn, err
		}
		conn = com.NewBufConn(ft, conn)
		if ft[0] == 0x16 && ft[1] == 0x03 {
			conn = tls.Server(conn, m.tlsconf)
		}
	}
	return conn, err

}

