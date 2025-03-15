package client

import (
	"github.com/sadeepa24/netshoot/result"
)

type clientNode interface {
	Test(host string) result.Result
	Start() error
	Close() error
}

