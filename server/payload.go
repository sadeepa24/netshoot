package server

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
)

type PayloadServer struct {
	ctx context.Context
	cancel context.CancelFunc

	payloadFirstSorted [][]byte //this is first part of all payload with sorted
	payloads           []com.Payload // all payloads sorted
	logger *zap.Logger

	listner net.Listener

	writeTimeout time.Duration
	readTimeout time.Duration

	bytepool sync.Pool
	mixedHandler *MixedHandler

}


func NewPayloadServer(ctx context.Context, logger *zap.Logger, conf config.PayloadServer) (*PayloadServer, error) {
	ps := &PayloadServer{
		logger: logger,
		writeTimeout: conf.WriteTimeout(),
		readTimeout: conf.ReadTimeout(),
	}
	ps.ctx, ps.cancel = context.WithCancel(ctx)
	allpl, err := com.ReadPayloadFile(conf.PayloadFile)
	if err != nil {
		return nil, err
	}
	ps.payloads = allpl
	ps.payloadFirstSorted = allpl.FirstPart()


	laddr, err := net.ResolveTCPAddr("tcp", conf.Ls.ListenAddr)
	if err != nil {
		return nil, err
	}
	ps.listner, err = net.ListenTCP("tcp", laddr)
	if err != nil {
		return nil, err
	}
	if conf.Ls.Tls.Enabled {
		ps.mixedHandler, err = NewMixedHandler(ctx, conf.Ls.Tls)
		if err != nil {
			return nil, err
		}
	}
	ps.bytepool = sync.Pool{
		New: func() any {
			return make([]byte, len(ps.payloadFirstSorted[len(ps.payloadFirstSorted)-1]))
		},
	}

	return ps, nil
}


func (p *PayloadServer) Start() error {
	p.logger.Debug("payload server started")
	go func() {
		for {
			if p.ctx.Err() != nil {
				break
			}
			conn, err := p.listner.Accept()
			if err != nil {
				continue
			}
			p.logger.Debug("connection recived " + conn.RemoteAddr().String())
			if tlsconn, ok := conn.(*tls.Conn); ok {
				p.logger.Debug("server tls handshake success", zap.String("sni", tlsconn.ConnectionState().ServerName))
			}
			go p.handleconn(conn)
		}
	}()

	return nil //satisfy the interface
	
}
func(p *PayloadServer) Close() error {
	p.listner.Close()
	p.cancel()
	return nil
}


func (p *PayloadServer) handleconn(conn net.Conn) {
	var (
		err error
		payloadNumber int
	)
	if p.mixedHandler != nil {
		conn, err = p.mixedHandler.handle(conn)
		if err != nil {
			p.logger.Error("Connection Prehandle Err MixedHandler err ", zap.Error(err))
			return
		}
	}
	conn.SetDeadline(time.Now().Add(p.readTimeout))
	defer conn.Close()
	payloadNumber, err = p.detectPayloadStrict(conn)
	if err != nil {
		p.logger.Error("Payload Detect Failed: ", zap.Error(err))
		return
	}
	p.logger.Debug("payload detected succesfully")
	
	payload := p.payloads[payloadNumber]
	host, err := payload.ReadAfterFirst(conn)
	if err != nil {
		p.logger.Error("Payload Read Error After Detect Payload: ", zap.Error(err), zap.String("host", host))
		return
	}
	p.logger.Info("recived host in payload", zap.String("host", host), zap.Int("payload", payloadNumber))
	conn.SetDeadline(time.Now().Add(p.writeTimeout))
	_, err = payload.WriteRes(conn)
	if err != nil {
		p.logger.Error("Payload Response Write Err: ", zap.Error(err))
		return
	}
	err = Speedtest(conn)
	if err != nil {
		p.logger.Error("Speedtest Failed: ", zap.Error(err))
	}
}

func (p *PayloadServer) detectPayloadStrict(conn net.Conn) (int, error) {
	payloadNumber := 0
	totalread := 0
	cache := p.bytepool.Get().([]byte)
	//lint:ignore SA6002 reason
	defer p.bytepool.Put(cache)
	for i, pload := range p.payloadFirstSorted {
		_, err := io.ReadFull(conn, cache[totalread:totalread+(len(pload)-totalread)])
		if err != nil {
			return payloadNumber, err
		}
		if bytes.Equal(pload, cache[:len(pload)]) {
			payloadNumber = i
			break
		}
		totalread = len(pload)
		if i == len(p.payloadFirstSorted)-1 {
			return payloadNumber, ErrPayloadNotfound
		}
	}
	return payloadNumber, nil
}

var ErrPayloadNotfound = errors.New("not valid payload")

func (p *PayloadServer) detectPayloadEasy(conn net.Conn) (int, error) {
	payloadNumber := 0
	for i, pload := range p.payloadFirstSorted {
		oldlen := 0 
		if i > 0 {
			oldlen = len(p.payloadFirstSorted[i-1])
		}
		tothis := make([]byte, len(pload)-oldlen)
		_, err := io.ReadFull(conn, tothis)
		if err != nil {
			return payloadNumber, err
		}
		if bytes.Equal(pload[oldlen:], tothis) {
			payloadNumber = i
		}
	}
	return payloadNumber, nil
}