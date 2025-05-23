package server

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
)

type MixedHandler struct {
	tlsconf *tls.Config
	tlstimeout time.Duration
	ctx context.Context
}

func NewMixedHandler(ctx context.Context, tlsconf config.TlsServer) (*MixedHandler, error) {
	mixhndler := &MixedHandler{
		ctx: ctx,
	}
	if tlsconf.Enabled {
		certs, err := tls.LoadX509KeyPair(tlsconf.Cert, tlsconf.Key)
		if err != nil {
			return nil, err
		}
		mixhndler.tlsconf = &tls.Config{
			Certificates: []tls.Certificate{certs},
			InsecureSkipVerify: true,
			MinVersion: tls.VersionTLS10,
			MaxVersion: tls.VersionTLS13,
		}
		mixhndler.tlstimeout = tlsconf.TlsTimeout()
	}
	return mixhndler, nil
}

func (m *MixedHandler) handle(conn net.Conn) (net.Conn, error) {
	var err error
	if m.tlsconf != nil {
		ft := make([]byte, 5)		
		conn.SetDeadline(time.Now().Add(m.tlstimeout))
		_, err := conn.Read(ft)
		if err != nil {
			conn.Close()
			return conn, err
		}
		conn = com.NewBufConn(ft, conn)
		if ft[0] == 0x16 && ft[1] == 0x03 {
			tlsconn := tls.Server(conn, m.tlsconf)
			timeoutctx, cancel := context.WithTimeout(m.ctx, m.tlstimeout)
			err = tlsconn.HandshakeContext(timeoutctx)
			cancel()
			if err != nil {
				tlsconn.Close()
				return nil, err
			}
			conn = tlsconn
		}
	}
	return conn, err
}