package result

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

/*
in progrees

checked host count
tcp failed count
other error
succesfull count

*/

const DefaultFlush = 3

type ProgressConf struct {
	CheckedHost int `json:"checked_host"`
	TotalTCPFail int `json:"total_tcp_fail"`
	LastHost string `json:"last_host"`
	TotalSuccess int `json:"total_success"`
}


type progress struct {
	conf ProgressConf
	tcpFailStreak int

	file *os.File

	tcpFailThreshHold int

	path string
	flushcounter  int
}

func newProgress(path string, tcpfailthresd int) (*progress, error) {
	if tcpfailthresd < 20 {
		tcpfailthresd = 20
	}
	pr := &progress{
		path: path,
		flushcounter: DefaultFlush,
		tcpFailThreshHold: tcpfailthresd,
	}
	return pr, nil
}

func (p *progress) Start() error {
	var err error
	p.file, err = os.OpenFile(p.path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	if err = p.loadconf(); err != nil {
		return err
	}
	return nil
}


func (p *progress) Close() error {
	p.flush()
	return p.file.Close()
}

// only required when starting so no race condition
func (p *progress) currentConf() ProgressConf {
	return p.conf
}


func (p *progress) Update(res []Result, order []string) bool {
	if p.conf.TotalTCPFail >= p.tcpFailThreshHold {
		return false
	}
	
	if len(order) == 0 {
		return true
	}
	p.conf.CheckedHost += len(order)
	for _, v := range res {
		p.conf.TotalTCPFail += v.TcpFailCount()
		p.conf.TotalSuccess += v.SuccesCount()
	}
	p.conf.LastHost = order[len(order)-1]
	p.flushcounter--
	if p.flushcounter == 0 {
		p.flush()
	}
	return true
}

func (p *progress) flush() {
	p.flushcounter = DefaultFlush
	p.file.Seek(0, 0)
	m, err := p.file.Write(
		[]byte( fmt.Sprintf(
			`{"checked_host": %d, "total_tcp_fail": %d, "last_host": "%s", "total_success": %d}`,
			p.conf.CheckedHost, p.conf.TotalTCPFail, p.conf.LastHost, p.conf.TotalSuccess,
		)  ),
	)
	if err == nil {
		p.file.Truncate(int64(m))
	}
}

func (p *progress) loadconf() error {
	file, err := io.ReadAll(p.file)
	if err != nil  {
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, &p.conf)
}