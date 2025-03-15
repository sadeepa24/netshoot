package hostmanager

import (
	"bytes"
	"errors"
	"io"
	"os"

	config "github.com/sadeepa24/netshoot/configs"
)


const DefaultRead = 1024 * 100 //(100kb)


type hostfile struct {
	buffer     []byte
	hostBuf []string
	file       *os.File
	path string
	totalread int
	initialized bool
	maxconcurrent int

	avbl bool
	fileover bool
	filesize int64
}

func (h *hostfile) Close() error {
	return h.file.Close()
}

func newfile(conf config.Hostfile) (*hostfile, error){
	fs := &hostfile{
		path: conf.Hostfile,
		maxconcurrent: conf.MaxConcurrent,

	}
	return fs, nil
}

func (h *hostfile) initialize(startfrom int) error {
	if h.initialized {
		return nil
	}
	var err error
	h.file, err = os.OpenFile(h.path, os.O_CREATE|os.O_RDWR, 0644)
	
	if err != nil {
		return err
	}
	fsStat, err := h.file.Stat()
	if err != nil {
		return err
	}

	if fsStat.Size() == 0 {
		h.file.Close()
		return errors.New("empty file")
	}
	h.filesize = fsStat.Size()
	h.avbl = true
	h.initialized = true

	return h.freeRead(startfrom)
}


func (h *hostfile) next() []string {
	if len(h.hostBuf) > h.maxconcurrent {
		next := h.hostBuf[:h.maxconcurrent]
		h.hostBuf = h.hostBuf[h.maxconcurrent:]
		return next
	}
	if h.fileover {
		h.avbl = false
		return h.hostBuf

	}
	h.fill()
	return h.next()
}

func (h *hostfile) nextone() (bool, string) {
	if len(h.hostBuf) > 0 {
		next := h.hostBuf[0]
		h.hostBuf = h.hostBuf[1:]
		return true, next
	}
	if h.fileover {
		return false, ""
	}
	h.fill()
	return h.nextone()
}

func (h *hostfile) freeRead(checkCount int) error {
	if checkCount < 0 {
		return errors.New("invalid check count")
	}
	if len(h.hostBuf) >= checkCount {
		h.hostBuf = h.hostBuf[checkCount:]
		return nil
	}
	if h.fileover {
		return errors.New("no more data: may be host file change")
	}

	oldLen := len(h.hostBuf)
	
	h.hostBuf = nil
	
	h.fill()
	return h.freeRead(checkCount-oldLen)
	
}


func (h *hostfile) fill() {
	h.read()
	s := 0
	for {
		s = bytes.IndexByte(h.buffer, '\n')
		if s == -1 {
			if h.fileover {
				h.hostBuf = append(h.hostBuf, string(removeCR(h.buffer)))
			}
			break
		}
		h.hostBuf = append(h.hostBuf, string(removeCR(h.buffer[:s])))
		if s+1 < len(h.buffer) {
			h.buffer = h.buffer[s+1:]
			continue
		}
		break
	}
	
}

func (h *hostfile) read() error {
	nextbufsize := DefaultRead
	if DefaultRead > int(h.filesize)-h.totalread {
		nextbufsize = int(h.filesize) - h.totalread
	}
	newbuf := make([]byte, nextbufsize+len(h.buffer))

	copy(newbuf, h.buffer)
	n, err := h.file.Read(newbuf[len(h.buffer):])
	h.buffer = newbuf
	h.totalread += n
	if h.totalread == int(h.filesize) {
		h.fileover = true
	}
	if err != nil {
		if errors.Is(err, io.EOF) {
			h.fileover = true
			return nil
		}
		return err
	}
	return err
}

func (h *hostfile) available() bool {
	return h.avbl
}

func removeCR(in []byte) []byte {
	if len(in) > 0 && in[len(in)-1] == '\r' {
		return in[0 : len(in)-1]
	}
	return in
} 