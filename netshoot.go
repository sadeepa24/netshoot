package netshoot

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/sadeepa24/netshoot/client"
	config "github.com/sadeepa24/netshoot/configs"
	"github.com/sadeepa24/netshoot/hostmanager"
	"github.com/sadeepa24/netshoot/result"
	"github.com/sadeepa24/netshoot/server"
	"go.uber.org/zap"
)

type Netshoot struct {
	ctx context.Context
	hostManager *hostmanager.HostManager
	client *client.Client
	server *server.Server
	result *result.ResultWriter

	logger *zap.Logger
	stop chan struct{}

	closed *atomic.Bool
}

func New(ctx context.Context, logger *zap.Logger, conf config.Config) (*Netshoot, error) {
	netshoot := &Netshoot{
		ctx: ctx,
		stop: make(chan struct{}),
		logger: logger,
		closed: new(atomic.Bool),
	}

	var err error
	netshoot.client, err = client.NewClient(ctx, conf.Client, logger)
	if err != nil {
		return nil, err
	}
	if netshoot.client.NodeCount() > 0 {
		if netshoot.result, err = result.NewResultWriter(ctx, conf.Result, netshoot.stop, logger); err != nil {
			return nil, errors.New("result write create failed: " + err.Error())
		}
		if netshoot.hostManager, err = hostmanager.New(ctx, netshoot.client, netshoot.result, logger, conf.Host, netshoot.stop); err != nil {
			return nil, errors.New("host manager create failed: " + err.Error())
		}

	} else {
		netshoot.client = nil
	}
	netshoot.server, err = server.NewServer(ctx, logger, conf.Server)
	if err != nil {
		return nil, err
	}
	if netshoot.server.NodeCount() == 0 {
		netshoot.server = nil
	}

	if netshoot.server == nil && netshoot.client == nil {
		return nil, errors.New("no client or server created, should recheck config")
	}
	return netshoot, nil
}

func (n *Netshoot) Start() error {
	var err error
	if n.server != nil {
		if err = n.server.Start(); err != nil {
			return err
		}
	}
	if n.client != nil {
		if err = n.client.Start(); err != nil{return err}
		if err = n.result.Start(); err != nil{return err}
		if err = n.hostManager.Start(); err != nil{return err}
	}
	go func() {
		select {
		case <- n.stop:
			n.logger.Debug("stop signal recived")
			err = n.Close()
			if err != nil {
				n.logger.Error(" netshoot closing err: ", zap.Error(err))
			}
			return
		case <- n.ctx.Done():
		}
	}()
	n.logger.Debug("Netshoot Started")
	return nil
}

func (n *Netshoot) Close() error {
	n.logger.Debug("closing Netshoot...")
	if n.closed.Swap(true) {
		return nil
	}
	if n.client != nil {
		n.hostManager.Close()
		n.logger.Debug("hostmanager closed")
		n.client.Close()
		n.logger.Debug("client closed")
		n.result.Close()
		n.logger.Debug("result closed")
	}
	if n.server != nil {
		n.server.Close()
	}
	return nil
}