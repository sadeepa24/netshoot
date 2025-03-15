package com

import (
	"sync"

	"github.com/sadeepa24/netshoot/result"
)

type Getresult struct {
	mu *sync.RWMutex
	wg *sync.WaitGroup
	currentRes []result.Result
}

func NewrsGet() *Getresult {
	return &Getresult{
		mu: &sync.RWMutex{},
		wg: &sync.WaitGroup{},
	}
}

func (r *Getresult) Reset(expectNext int) {
	r.wg.Add(expectNext)
	r.mu.Lock()
	r.currentRes = []result.Result{}
	r.mu.Unlock()
	
}

func (r *Getresult) UploadResult(res []result.Result) {
	r.wg.Done()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.currentRes = append(r.currentRes, res...)
}

func (r *Getresult) Wait() []result.Result {
	r.wg.Wait()
	return r.currentRes
}