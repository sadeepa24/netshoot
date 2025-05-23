package hostmanager

import (
	"context"
	"errors"
	"time"

	"github.com/sadeepa24/netshoot/client"
	com "github.com/sadeepa24/netshoot/common"
	config "github.com/sadeepa24/netshoot/configs"
	"github.com/sadeepa24/netshoot/result"
	"go.uber.org/zap"
)

type HostManager struct {
	
	ctx context.Context
	
	logger *zap.Logger
	file *hostfile
	resultRW *result.ResultWriter
	client *client.Client
	resultgt *com.Getresult

	done chan struct{}
	stopsig chan struct{}

	interval time.Duration
	// over *atomic.Bool
	

}

func New(ctx context.Context, client *client.Client, resultRW *result.ResultWriter, logger *zap.Logger, conf config.HostMgConf, stopsig chan struct{}) (*HostManager, error) {
	
	hostmg := &HostManager{
		ctx: ctx,
		done: make(chan struct{}) ,
		resultgt: com.NewrsGet(),
		client: client,
		resultRW: resultRW,
		logger: logger,
		stopsig: stopsig,
		interval: conf.Interval,
		//over: new(atomic.Bool),
	}

	var err error
	hostmg.file, err = newfile(conf.Hostfile)
	if err != nil {
		return nil, err
	}
	
	return hostmg, nil
}

func (r *HostManager) Start() error { 
	err := r.file.initialize(r.resultRW.Progres().CheckedHost)
	if err != nil {
		return err
	}
	go r.run()
	return nil
}

func (r *HostManager) Close() error {
	
	select {
	case <-r.done:
	default:
		close(r.done)
	}

	return 	r.file.Close()
}




func (r *HostManager)  run() {	
	var nextCheck []string
	for r.file.available() {
		//gracefull shutdonw stuff
		select {
			case <- r.done:
				//TODO: add gracefull shutdown to here
				r.logger.Debug("done signal recived closing hostmanager run loop")
				return
			case <- r.ctx.Done():
				//forceclosing
				r.logger.Warn("force context canceled closing run loop")
				return
			default:
			
		}
		if r.interval > 0 {
			time.Sleep(r.interval)
		}
		nextCheck = r.file.next()
		r.resultgt.Reset(len(nextCheck))
		for _, host := range nextCheck {
			go r.client.MakeTest(host, r.resultgt)
		}
		r.resultRW.Write(r.resultgt.Wait(), nextCheck)
	}
	close(r.done)
	r.stopsig <- struct{}{}
}



var (
	ErrForceCancel = errors.New("force canceld run loop ")
)