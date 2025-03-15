package config

import (
	"crypto/tls"
	"net"
	"time"
)


var (
	DefaultDTimeout = 300 * time.Millisecond // Dialer Default Timeout
	DefaultTAuthTimeout = 300 * time.Millisecond // tls auth timeout
	
	//client
	DefaultClientWTimeout = 300 * time.Millisecond // Default Payload Write Timeout
	DefaultClientRTimeout = 300 * time.Millisecond // Default Payload Response Read Timeout
	
	//server
	DefaultSrvWTimeout = 300 * time.Millisecond // Default Payload Write Timeout
	DefaultSrvtRTimeout = 300 * time.Millisecond // Default Payload Response Read Timeout
)


type HostMgConf struct {
	Hostfile
}

type PayloadSender struct {
	
	DialerTimeout string  `json:"dialer_timeout"`
	HandshakeRetry int  `json:"handshake_retry"`
	ServerAddr string  `json:"server_addr"`
	Tls TlsConf  `json:"tls"`
	PayloadFile string  `json:"payload_file"`
	TestBufSize int  `json:"speedtest_size"`
	Interface string  `json:"net_iface"`
	Local_addr string   `json:"local_addr"`

	RTimeout string  `json:"read_timeout"`
	WTimeout string  `json:"write_timeout"`
}

func (p PayloadSender) DTimeout() time.Duration {
	if p.DialerTimeout == "" {
		return DefaultDTimeout
	}
	time, err := time.ParseDuration(p.DialerTimeout)
	if err != nil {
		return DefaultDTimeout
	}
	return time
}

func (p PayloadSender) LocalAddr() net.Addr {
	if p.Local_addr != "" {
		addr, err := net.ResolveTCPAddr("ip", p.Local_addr+ ":0")
		if err == nil {
			return addr
		}
	}
	
	if p.Interface != "" {
		iface, err := net.InterfaceByName(p.Interface)
		if err == nil {
			aadres, err := iface.Addrs()
			if err == nil {
				for _, addr := range aadres  {
					if ad, ok := addr.(*net.IPNet); ok && ad.IP.To4() != nil {
						return &net.TCPAddr{
							IP:   ad.IP,
							Port: 0,
							Zone: iface.Name,
						}
					}
				}
			}
		}
	}
	return nil
}
func (p PayloadSender) SpeedTestBuf() int {
	//TODO: parse p.Testbufsize and convert it to bytesize
	return p.TestBufSize
	//return p.TestBufSize
}

func (p PayloadSender) TlsAuthTimeout()time.Duration  {
	if p.Tls.AuthTimeout == "" {
		return DefaultTAuthTimeout
	}
	time, err := time.ParseDuration(p.Tls.AuthTimeout)
	if err != nil {
		return DefaultTAuthTimeout
	}
	return time
}

func (p PayloadSender) WriteTimeout()time.Duration  {
	if p.WTimeout == "" {
		return DefaultClientWTimeout
	}
	time, err := time.ParseDuration( p.WTimeout)
	if err != nil {
		return DefaultClientWTimeout
	}
	return time
}

func (p PayloadSender) ReadTimeout()time.Duration  {
	if p.RTimeout == "" {
		return DefaultClientRTimeout
	}
	time, err := time.ParseDuration(p.RTimeout)
	if err != nil {
		return DefaultClientRTimeout
	}
	return time
}

type TlsConf struct {
	Enabled bool  `json:"enabled"`
	Min string  `json:"min_version"`
	Max string  `json:"max_version"`
	AuthTimeout string  `json:"auth_timeout"`
	Insecure bool  `json:"insecure"`
	NextProt []string  `json:"next_proto"`


}

func (t *TlsConf) Maxversion() uint16 {
	return tls.VersionTLS12
}
func (t *TlsConf) Minversion() uint16 {
	return tls.VersionTLS11
}


type ClientNode struct {
	Type string  `json:"type"`
	Disabled bool `json:"disabled"`
	*PayloadSender

}


type Result struct {
	OutputFile      string `json:"output_file"`
	ProgressFile    string `json:"progress_file"`
	TCPFailThreshHold int   `json:"tcp_fail_threshold"`
}


type Hostfile struct {
	MaxConcurrent int  `json:"max_concurrent"`
	Hostfile string  `json:"host_file"`
}


//client side

type Client struct {
	Nodes []ClientNode  `json:"nodes"`
}








// server side
type Server struct {
	Nodes []ServerNode  `json:"nodes"`
}

type ServerNode struct {
	Type string  `json:"type"`
	Disabled bool `json:"disabled"`
	*PayloadServer
}

type PayloadServer struct {
	PayloadFile string  `json:"payload_file"`
	Ls LsConfig  `json:"listen_conf"`
	RTimeout string  `json:"read_timeout"`
	WTimeout string  `json:"write_timeout"`
}

func (p PayloadServer) WriteTimeout()time.Duration  {
	if p.WTimeout == "" {
		return DefaultSrvWTimeout
	}
	time, err := time.ParseDuration( p.WTimeout)
	if err != nil {
		return DefaultSrvWTimeout
	}
	return time
}

func (p PayloadServer) ReadTimeout()time.Duration  {
	if p.RTimeout == "" {
		return DefaultSrvtRTimeout
	}
	time, err := time.ParseDuration(p.RTimeout)
	if err != nil {
		return DefaultSrvtRTimeout
	}
	return time
}

type LsConfig struct {
	ListenAddr string  `json:"listen"`
	Tls TlsServer  `json:"tls"`
}

type TlsServer struct {
	Enabled bool  `json:"enabled"`
	Cert string  `json:"cert"`
	Key string  `json:"key"`

}
