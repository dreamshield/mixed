package pool

import (
	"errors"
	"fmt"
)

// Goroutine Pool
type RoutinePool interface {
	Take()
	Return()
	Active() bool
	Total() uint
	Remainder() uint
}

//  worker pool
type WorkerPool struct {
	total    uint
	ticketCh chan byte
	active   bool
}

func NewWorkerPool(total uint) (RoutinePool, error) {
	gt := WorkerPool{}
	if !gt.init(total) {
		errMsg :=
			fmt.Sprintf("The goroutine worker pool can NOT be initialized! (total=%d)\n", total)
		return nil, errors.New(errMsg)
	}
	return &gt, nil
}

func (gt *WorkerPool) init(total uint) bool {
	if gt.active {
		return false
	}
	if total == 0 {
		return false
	}
	ch := make(chan byte, total)
	n := int(total)
	for i := 0; i < n; i++ {
		ch <- 1
	}
	gt.ticketCh = ch
	gt.total = total
	gt.active = true
	return true
}

func (gt *WorkerPool) Take() {
	<-gt.ticketCh
}

func (gt *WorkerPool) Return() {
	gt.ticketCh <- 1
}

func (gt *WorkerPool) Active() bool {
	return gt.active
}

func (gt *WorkerPool) Total() uint {
	return gt.total
}

func (gt *WorkerPool) Remainder() uint {
	return uint(len(gt.ticketCh))
}
