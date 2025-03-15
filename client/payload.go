package client

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"time"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
	"github.com/sadeepa24/netshoot/result"
	"go.uber.org/zap"
)

type PayloadSender struct {

	ctx context.Context

	Dialer net.Dialer
	payloads []com.Payload
	handshakeretry int
	serveraddr string
	tlsconf *tls.Config

	tlstimeout time.Duration
	speedtestSize int

	writeTimeout time.Duration
	readTimeout time.Duration

	logger *zap.Logger
}


func NewPayloadSender(ctx context.Context, conf config.PayloadSender, logger *zap.Logger) (*PayloadSender, error){
	sender := PayloadSender{
		Dialer: net.Dialer{
			Timeout: conf.DTimeout(),
			LocalAddr: conf.LocalAddr(),
		},
		serveraddr: conf.ServerAddr,
		speedtestSize: conf.SpeedTestBuf(),
		tlstimeout: conf.TlsAuthTimeout(),

		writeTimeout: conf.WriteTimeout(),
		readTimeout: conf.ReadTimeout(),
		handshakeretry: conf.HandshakeRetry, //TODO: change later
		logger: logger,
		ctx: ctx,
	}
	var err error
	if sender.payloads, err = com.ReadPayloadFile(conf.PayloadFile); err != nil {
		return nil, err
	}
	if conf.Tls.Enabled {
		sender.tlsconf = &tls.Config{
			InsecureSkipVerify: conf.Tls.Insecure,
			NextProtos: conf.Tls.NextProt,
			MaxVersion: conf.Tls.Maxversion(),
			MinVersion: conf.Tls.Minversion(),	
		}
	}

	return &sender, nil
}

func (p *PayloadSender) Test(host string) result.Result {
	res := result.PayloadResult{
		ComResult: result.ComResult{
			Host: host,
		},
	}
	for _, payload := range p.payloads {
		sres := result.SinglePayload{
			PayloadName: payload.Name(),
		}
		conn, err := p.Dial()
		if err != nil {
			sres.TcpFailed = true
			sres.Error = err.Error()
			res.PayloadInfo = append(res.PayloadInfo, sres)
			continue
		}
		p.logger.Debug("dialing success host: " + host)

		if p.tlsconf != nil {
			
			tlsconn := tls.Client(conn, &tls.Config{  //we have to copy tls.config to use concurrently because servername changes everytime
				ServerName: host,
				MaxVersion: p.tlsconf.MaxVersion,
				MinVersion: p.tlsconf.MinVersion,
				NextProtos: p.tlsconf.NextProtos,
				InsecureSkipVerify: p.tlsconf.InsecureSkipVerify,
			})
			sres.Tls = result.TlsInfo{
				Servername: host,
			}
			//TODO: fix later
			timeoutctx, cancel := context.WithTimeout(p.ctx, p.tlstimeout)
			err = tlsconn.HandshakeContext(timeoutctx)
			cancel()
			if err != nil {
				sres.Tls.Failed = true
				sres.Tls.Error = err.Error()
				res.PayloadInfo = append(res.PayloadInfo, sres)
				conn.Close()
				continue
			}
			conn = tlsconn
		}
		conn.SetDeadline(time.Now().Add(p.writeTimeout))
		if err = payload.WriteTo(conn, host); err != nil {
			sres.Error = "payload write err: " + err.Error()
			res.PayloadInfo = append(res.PayloadInfo, sres)
			conn.Close()
			continue
		}
		p.logger.Debug("payload write success host: " + host)
		conn.SetDeadline(time.Now().Add(p.readTimeout))
		if _, err = payload.ReadRes(conn); err != nil {
			sres.Error = "payload response read err: " + err.Error()
			res.PayloadInfo = append(res.PayloadInfo, sres)
			conn.Close()
			continue
		}
		p.logger.Debug("payload read success host: "+ host)
		sres.Maybe = true
		sres.MaxSpeed, err = speedtest(conn, p.speedtestSize)
		if err != nil {
			sres.Error = "speedtest err: " + err.Error()
			res.PayloadInfo = append(res.PayloadInfo, sres)
			conn.Close()
			continue
		}
		p.logger.Debug("Speedtest success", zap.Float64("speed", sres.MaxSpeed))
		sres.Success = true
		res.PayloadInfo = append(res.PayloadInfo, sres)
		conn.Close()

	}
	res.PreProcess()
	return res

}

func (p *PayloadSender) Dial() (net.Conn, error) {
	var (
		conn net.Conn
		err error
	)
	for i := 0; i < p.handshakeretry; i++ {
		conn, err = p.Dialer.Dial("tcp",  p.serveraddr)
		if err == nil {
			return conn, nil
		}
	}
	return nil, errors.Join(ErrAllRetryFailed, err,)
}



var (
	ErrAllRetryFailed  = errors.New("all tcp retry failed")
)

//satisfy the interface
func (p *PayloadSender) Start() error {
	return nil
}
func (p *PayloadSender) Close() error {
	return nil
}