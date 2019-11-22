package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"
)

type ProcessStatus struct {
	sync.RWMutex
	rusage      *syscall.Rusage
	lastCpuTime int64
	stopSignal  chan bool
	stopped     bool

	CpuUsage float64
	MemRss   int64
}

func NewProcessStatus() *ProcessStatus {
	p := new(ProcessStatus)
	p.rusage = new(syscall.Rusage)

	go func() {
		timer := time.Tick(time.Second)
		for {
			select {
			case <-timer:
				p.Update()
			case <-p.stopSignal:
				return
			}
		}
	}()

	return p
}

func (p *ProcessStatus) Update() {
	e := syscall.Getrusage(syscall.RUSAGE_SELF, p.rusage)
	if e != nil {
		fmt.Println("failed-to-get-rusage", e)
	}

	p.Lock()
	defer p.Unlock()
	p.MemRss = int64(p.rusage.Maxrss)

	t := p.rusage.Utime.Nano() + p.rusage.Stime.Nano()
	p.CpuUsage = float64(t-p.lastCpuTime) / float64(time.Second.Nanoseconds())
	p.lastCpuTime = t
}

func (p *ProcessStatus) StopUpdate() {
	p.Lock()
	defer p.Unlock()
	if !p.stopped {
		p.stopped = true
		p.stopSignal <- true
		p.stopSignal = nil
	}
}
