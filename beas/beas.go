// BEAS is a data dealing progress based on master-worker model
// You can modify the model packge to satisfy your own business
package main

import (
	"fmt"
	"mixed/beas/file"
	"mixed/beas/model"
	"mixed/beas/pool"
	"mixed/beas/setting"
	"sync"
)

var workers pool.RoutinePool
var err error
var wg sync.WaitGroup

func init() {
	// parse parameters
	setting.ParseParams()
	// initialize DB
	model.InitDBEngine()
	// initialize worker
	workers, err = pool.NewWorkerPool(setting.WorkerNum)
	if err != nil {
		panic("Init worker pool fail msg :" + err.Error())
	}
}

func main() {

	if setting.ShowVersion {
		setting.PrintVersion()
	}

	if setting.ShowHelp {
		setting.Usage()
	}
	// master-worker control
LOOP:
	for {
		if setting.NowPage > setting.TotalPage {
			break LOOP
		}
		format := "curpage = %v, begin = %v, end=%v\n"
		file.WriteDataFile(file.LogFd, format, setting.NowPage, setting.ProbeBegin, setting.ProbeEnd)
		workers.Take()
		go model.GetOrderWorker(setting.ProbeBegin, setting.ProbeEnd, workers, wg)
		setting.NowPage++
		setting.ProbeBegin = setting.ProbeEnd + 1
		setting.ProbeEnd = setting.ProbeEnd + setting.LimitPerTime
	}
	fmt.Println("Order Product Search Done")
	wg.Wait()
	// TODO: notify to close the file handler
	fmt.Println("Generate Updata SQL Done")
}
