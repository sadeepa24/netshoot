package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/sadeepa24/netshoot"
	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	/*
	ss := config.Config{
		Client: config.Client{
			Nodes: []config.ClientNode{
				{
					Type: "payload",
					PayloadSender: &config.PayloadSender{
						HandshakeRetry: 5,
						PayloadFile:    "payload.dt",
						TestBufSize:    3,
						ServerAddr:     "127.0.0.1:90",
						Tls: config.TlsConf{
							Enabled: false,
						},
					},
				},
			},
		},
		Host: config.HostMgConf{
			Hostfile: config.Hostfile{
				MaxConcurrent: 10,
				Hostfile:      "host.txt",
			},
		},
		Result: config.Result{
			OutputFile:   "out.json",
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
	}

	ssss, err := json.Marshal(ss)
	if err != nil {
		logger.Fatal(err.Error())
	}
	f, err :=os.OpenFile("config.json", os.O_CREATE|os.O_RDWR, 0644) 
	if err != nil {
		return
	}
	f.Write(ssss)
	return
	*/

	configb, err := os.ReadFile("config.json")
	if err != nil {
		logger.Fatal("config file open err: " + err.Error())
	}
	var config config.Config
	err = json.Unmarshal(configb, &config)
	if err != nil {
		logger.Fatal("config Unmarshal Err: " + err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	shoot, err := netshoot.New(ctx, logger, config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = shoot.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	shoot.Close()
	cancel()


}