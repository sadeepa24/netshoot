package client

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
	C "github.com/sadeepa24/netshoot/constant"
	"github.com/sadeepa24/netshoot/result"
	"go.uber.org/zap"
)

type Client struct {
	nodes []clientNode
}

func NewClient(ctx context.Context, conf config.Client, logger *zap.Logger) (*Client, error) {
	cl := &Client{}
	for i :=  range conf.Nodes {
		if conf.Nodes[i].Disabled {
			continue
		}
		nd, err := NewNode(ctx,  conf.Nodes[i], logger)
		if err != nil {
			return nil, errors.New("node create failed: " + err.Error())
		}
		cl.nodes = append(cl.nodes, nd)
	}
	
	return cl, nil
}

// after adding other node type they may have start, close
// but for now thease 2 do not need
func (c *Client) Start() error {
	for _, nd := range c.nodes {
		nd.Start()
	}
	return nil
}
func (c *Client) Close() error {
	for _, nd := range c.nodes {
		nd.Close()
	}
	return nil
}




func (c *Client) NodeCount() int {
	return len(c.nodes)
}

func (c *Client) MakeTest(host string, rg *com.Getresult) {
	allres := []result.Result{}
	for i := range c.nodes {
		allres = append(allres, c.nodes[i].Test(host))
	}
	rg.UploadResult(allres)
}

// bufsize should be in MB format 
// Ex - if bufsize 2MB bufsize = 2
// also Return Speed As MBps 
func speedtest(conn net.Conn, bufsize int) (float64, error) {
	conn.SetDeadline(time.Now().Add(20 * time.Duration(bufsize) * time.Second)) // expect minimum 50KBps
	ss := uint16(bufsize)
	err := binary.Write(conn, binary.LittleEndian, &ss)
	if err != nil {
		return 0, err
	}

	ch := make([]byte, 4)

	_, err = io.ReadFull(conn, ch)
	if err != nil {
		return 0, err
	}

	if string(ch) != "done" {
		return 0, ErrSpeedTestRes
	}

	mustread := bufsize * 1024 * 1024
	cache := make([]byte, 1024 * 1024)

	if len(cache) > mustread {
		cache = cache[:mustread]
	}

	n := 0
	st := time.Now()
	for {
		n, err = conn.Read(cache[:min(mustread, len(cache))])
		mustread -= n
		if mustread == 0 {
			if errors.Is(err, io.EOF) { // this happens due to server may close connection after writing all data
				err = nil
			}
			break
		}
		if mustread < 0 {
			err = ErrSpeedTestMalform
			break
		}
		if err != nil {
			return 0, err
		}

	}
	elp := time.Since(st)

	if elp.Seconds() == 0 {
		return -1, ErrSpeedTestNotime 
	}

	return float64(bufsize)/elp.Seconds(), err
}

var ErrSpeedTestRes = errors.New("speedtest response not recived")
var ErrSpeedTestMalform = errors.New("speedtest may be malformed and wrong")
var ErrSpeedTestNotime = errors.New("no time elpsed speedtest")


func NewNode(ctx context.Context, conf config.ClientNode, logger *zap.Logger) (clientNode, error) {
	switch conf.Type {
	case C.ClNodePayload:
		return NewPayloadSender(ctx, *conf.PayloadSender, logger)
	case C.ClNodeHttp:
		return nil, errors.New("no available yet")
	default:
		return nil, errors.New("Unknown Type: " + conf.Type)
	}

}