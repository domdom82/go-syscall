package main

import (
	"fmt"
	"math"
	"runtime"
	"time"
)

func main() {

	numCPU := runtime.NumCPU()

	// catch end signal
	end := make(chan bool, 1)

	// first start a processStatus
	pidStat := NewProcessStatus()

	// print ps status every 2 seconds
	go func() {
		tick := time.Tick(2 * time.Second)
		for {
			select {
			case <-end:
				break
			case <-tick:
				fmt.Printf("Raw CPU: %f CPU(%%): %.2f%%  Total CPU(%%): %.2f%%\n", pidStat.CpuUsage, pidStat.CpuUsage*100, pidStat.CpuUsage*100/float64(numCPU))
			}

		}
	}()

	// do some work
	for i := 0; i < numCPU; i++ {
		go func() {
			angle := 0.0
			for {
				_ = math.Sin(angle)
				_ = math.Cos(angle)
				_ = math.Tan(angle)

				angle += 0.5
				if angle > 360.0 {
					angle = 0.0
				}
			}
		}()
	}

	time.Sleep(30 * time.Second)
	end <- true
}
