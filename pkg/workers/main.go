package workers

import (
	"TestTask-events/pkg/app"
	"fmt"
	"github.com/pkg/errors"
	"runtime/debug"
	"sync"
)

type WorkerFunc func(val interface{}, resChan chan<- interface{}) error

type WorkerGroup struct {
	countWorkers int
	bufferSize   int
	wg           *sync.WaitGroup
	isStarted    bool
	ChIn         chan interface{}
	ChOut        chan interface{}
	core         app.Core
	workerFunc   WorkerFunc
}

func New(core app.Core, countWorkers int, bufferSize int, workerFunc WorkerFunc) *WorkerGroup {
	return &WorkerGroup{
		wg:           &sync.WaitGroup{},
		countWorkers: countWorkers,
		bufferSize:   bufferSize,
		ChIn:         make(chan interface{}, bufferSize),
		ChOut:        make(chan interface{}, bufferSize),
		core:         core,
		workerFunc:   workerFunc,
	}
}

func (w *WorkerGroup) Start() {
	w.isStarted = true
	w.wg.Add(w.countWorkers)
	for i := 0; i < w.countWorkers; i++ {
		go w.handle()
	}
	return
}

func (w *WorkerGroup) Stop() {
	w.wg.Wait()
	close(w.ChOut)
}

func (w *WorkerGroup) Send(val interface{}) error {
	w.ChIn <- val
	return nil
}

func (w *WorkerGroup) SetInputChan(ch chan interface{}) error {
	if w.isStarted {
		return fmt.Errorf("failed set input chan: workers already started")
	}
	w.ChIn = ch
	return nil
}

func (w *WorkerGroup) SetResChan(ch chan interface{}) error {
	if w.isStarted {
		return fmt.Errorf("failed set input chan: workers already started")
	}
	w.ChIn = ch
	return nil
}

func (w *WorkerGroup) handle() {
	defer w.wg.Done()

	for val := range w.ChIn {
		w.do(val)
	}
}

func (w *WorkerGroup) do(val interface{}) {

	defer func() {
		if r := recover(); r != nil {
			w.core.Logger().Printf("(PANIC) %v %v", r, string(debug.Stack()))
		}
	}()

	err := w.workerFunc(val, w.ChOut)
	if err != nil {
		w.core.Logger().Println("(ERROR) " + errors.Wrap(err, "failed work").Error())
	}

}
