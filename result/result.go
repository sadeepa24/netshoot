package result

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
)

const writeChanBuf int = 200

type ResultWriter struct {
	
	ctx context.Context
	
	resultFile *os.File
	prog *progress
	path string

	writer resultWriter

	writeChan chan Result
	done,stopsig chan struct{}
	logger *zap.Logger

	closed bool

	mu sync.Mutex
}

func NewResultWriter(ctx context.Context, conf config.Result, signal chan struct{}, logger *zap.Logger) (*ResultWriter, error) {
	res := &ResultWriter{
		path: conf.OutputFile,
		ctx: ctx,
		logger: logger,
		done: make(chan struct{}),
		writeChan: make(chan Result, writeChanBuf),
		mu: sync.Mutex{},
		stopsig: signal,
	}
	var err error
	
	res.prog, err = newProgress(conf.ProgressFile, conf.TCPFailThreshHold)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ResultWriter) Start() error {
	err := r.prog.Start()
	if err != nil {
		return errors.New("progress start failed: " + err.Error())
	}
	r.resultFile, err = os.OpenFile(r.path, os.O_CREATE|os.O_RDWR, 0644) 
	if err != nil {
		return err
	}
	r.writer = resultWriter{
		internal: r.resultFile,
	}
	err = r.writer.initialize()
	if err != nil {
		return errors.New("result writer initilize failed: " + err.Error())
	}
	go r.writeloop()

	return nil
}



func (r *ResultWriter) Close() error {
	r.mu.Lock()
	r.closed = true
	r.writeChan <- nil
	r.mu.Unlock()
	<- r.done
	r.prog.Close()
	close(r.writeChan)
	close(r.done)
	return r.writer.Close()
}

func (r *ResultWriter) Write(in []Result, hostorder []string) {
	//TODO: process progress updates then upload result to chan according to it 
	r.mu.Lock()
	if r.closed {
		return 
	}
	r.mu.Unlock()
	do := r.prog.Update(in, hostorder) 
	r.mu.Lock()
	if !do {
		r.logger.Info("program stopping due to exceeding TCP fail threshold")
		r.stopsig <- struct{}{}
		r.closed = true
	}
	for _, res := range in {
		if res == nil {
			continue
		}
		r.writeChan <- res
	}
	r.mu.Unlock()
}

func(r *ResultWriter)  writeloop() {
	for res := range r.writeChan {		
		if res == nil {
			break
		}
		if (r.ctx.Err() != nil) && len(r.writeChan) == 0 {
			r.logger.Debug("context canceled or closed result_writer, closing writer loop")
			return
		}
		resb, err := res.MarshalJSON()
		if err != nil {
			r.logger.Error("result marshaling failed for " + res.GetHost(), zap.Error(err))
			continue
		}
		_, err = r.writer.Write(resb)
		if err != nil {
			r.logger.Error("result writing failed for "+ res.GetHost(), zap.Error(err))
			continue
		}
	}
	r.done <- struct{}{}
}

func (r *ResultWriter) Progres() ProgressConf {
	return r.prog.currentConf()
}



type resultWriter struct {
	internal io.ReadWriter
	initiated bool
}

func (r *resultWriter) Write(in []byte) (int, error) {
	if !r.initiated {
		if err := r.initialize(); err != nil {
			return 0, err
		}
	}
	n, err := r.internal.Write(in)
	r.internal.Write([]byte(","))
	return n, err
}

func (r *resultWriter) Close() error {
	if !r.initiated {
		return nil
	}
	if file, ok := r.internal.(*os.File); ok {
		file.Seek(-1, io.SeekEnd)
		file.Write([]byte("]"))
		return file.Close()
	}
	return nil
}



func (r *resultWriter) initialize() error {
	if file, ok := r.internal.(*os.File); ok {
		fsInfo, err := file.Stat()
		if err != nil { return err}
		if fsInfo.Size() == 0 {
			r.initiated = true
			r.internal.Write([]byte("["))
		} else {
			_, err = file.Seek(-1, io.SeekEnd)
			if err != nil {
				return errors.New("failed to seek to end of file: " + err.Error())
			}
			verify := make([]byte, 1)
			_, err = r.internal.Read(verify)
			if err != nil { return err }
			
			
			if verify[0] == '}'  {
				_, err = file.Seek(-1, io.SeekEnd)
				if err != nil {
					return errors.New("failed to seek to end of file: " + err.Error())
				}
				r.internal.Write([]byte(","))

			} else if verify[0] != ','  {
				readpos := fsInfo.Size() - 1
				for {
					_, err := file.ReadAt(verify, readpos)
					readpos--
					if readpos < 0 {
						return errors.New("malformed output file")
					}
					if err != nil { return err }

					if verify[0] == ',' {
						file.Seek(readpos+2, io.SeekStart)
						break
					}
					if verify[0] == '}' {
						readpos++
					}
					if verify[0] == ']' ||  verify[0] == '}' {
						file.WriteAt([]byte(","), readpos+1)
						file.Seek(readpos+2, io.SeekStart)
						break
					} 
				}
			}
			r.initiated = true
		}
	}
	return nil
}

