package server

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"time"

	config "github.com/sadeepa24/netshoot/configs"
	C "github.com/sadeepa24/netshoot/constant"
	"go.uber.org/zap"
)


type serverNode interface {
	Start() error
	Close() error
}

type Server struct {
	nodes []serverNode
}

func NewServer(ctx context.Context, logger *zap.Logger, conf config.Server) (*Server, error) {
	srv := &Server{}
	
	for _, nd := range conf.Nodes {
		if nd.Disabled {
			continue
		}
		node, err := NewNode(ctx, logger, nd)
		if err != nil {
			return nil, err
		}
		srv.nodes = append(srv.nodes, node)
	}
	return srv, nil
}

func (c *Server) Start() error {
	for i := range c.nodes {
		if err := c.nodes[i].Start(); err != nil {
			return err
		}
	}
	return nil
}
func (c *Server) Close() error {
	for _, nd := range c.nodes {
		nd.Close()
	} // most time server close error deoes not matter 😅
	return nil
}


func (c *Server) NodeCount() int {
	return len(c.nodes)
}


func Speedtest(conn net.Conn) error {
	var size uint16
	err := binary.Read(conn, binary.LittleEndian, &size)
	
	if err != nil {
		return err
	}
	conn.SetDeadline(time.Now().Add(time.Duration(size) * 20 * time.Second)) // expect client have at least 50KBps connection 
	_,err = conn.Write([]byte("done"))
	if err != nil {
		return err
	}

	wr := make([]byte, 1024 * 1024 * 1)
	for i := range wr {
	    wr[i] = 0xFF  
	}
	mustwrite := int(size) * 1024 *1024
	
	if len(wr) > mustwrite {
		wr = wr[:mustwrite]
	}
	
	n := 0
	for {
		n, err = conn.Write(wr[:min(mustwrite, len(wr))])
		mustwrite -= n
		if mustwrite <= 0 {
			return nil
		}
		if err != nil{
			return err
		}
	}
}


func NewNode(ctx context.Context, logger *zap.Logger, conf config.ServerNode) (serverNode, error) {
	
	switch conf.Type {
	case C.ClNodePayload:
		return NewPayloadServer(ctx, logger, *conf.PayloadServer)
	case C.ClNodeHttp:
		return nil, errors.New("this feature is not developed yet, it may be available in the future")
	default:
		return nil, errors.New("Unknown Server Type: " + conf.Type)
	}

}