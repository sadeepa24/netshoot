package netshoot_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sadeepa24/netshoot"
	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
)

func TestNetshoot(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	logger = zap.NewNop()
	
	shoot, err := netshoot.New(context.Background(), logger, config.Config{
		Client: config.Client{
			Nodes: []config.ClientNode{
				{
					Type: "payload",
					PayloadSender: &config.PayloadSender{
						HandshakeRetry: 3,
						PayloadFile: "payload.dt",
						TestBufSize: 4,
						ServerAddr: "127.0.0.1:90",
						Tls: config.TlsConf{
							Enabled: false,
						},

					},
				},
			},
		},
		Host: config.HostMgConf{
			Hostfile: config.Hostfile{
				MaxConcurrent: 1,
				Hostfile: "host.txt",
			},
		},
		Result: config.Result{
			OutputFile: "out.json",
			ProgressFile: "prog.json",
		},
		Server: config.Server{
			Nodes: []config.ServerNode{
				{
					Type: "payload",
					PayloadServer: &config.PayloadServer{
						PayloadFile: "payload.dt",
						Ls: config.LsConfig{
							ListenAddr: "127.0.0.1:90",
							Tls: config.TlsServer{
								Enabled: false,
							},
						},
					},
				},
			},
		},

	})
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = shoot.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}
	fmt.Println("started")
	time.Sleep(1 * time.Second)

	fmt.Println("should stop")
	shoot.Close()
}