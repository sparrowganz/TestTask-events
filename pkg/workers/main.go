package workers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sparrowganz/TestTask-events/pkg/app"
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

func New(wg *sync.WaitGroup, core app.Core, countWorkers int, bufferSize int, workerFunc WorkerFunc) *WorkerGroup {
	return &WorkerGroup{
		wg:           wg,
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
		go func() {
			w.handle()
		}()
	}
	return
}

func (w *WorkerGroup) Stop() {
	close(w.ChIn)
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
		err := w.workerFunc(val, w.ChOut)
		if err != nil {
			w.core.Logger().Println("(ERROR) " + errors.Wrap(err, "failed work").Error())
		}
	}
	close(w.ChOut)
}