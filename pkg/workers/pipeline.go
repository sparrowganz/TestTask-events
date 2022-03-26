package workers

import (
	"fmt"
	"sync"
)

type Pipeline interface {
	Send(val interface{})
	GetResChan() chan interface{}
	Start()
	Stop()
}

type pipeData struct {
	wg      *sync.WaitGroup
	workers []*WorkerGroup
}

func NewPipeline(globalWg *sync.WaitGroup, workers ...*WorkerGroup) (Pipeline, error) {

	if len(workers) == 0 {
		return nil, fmt.Errorf("failed create pipeline: workers not found")
	}

	pipeline := &pipeData{
		wg:      globalWg,
		workers: workers,
	}

	for idx := range pipeline.workers {
		//Set current outChan
		if idx < len(pipeline.workers)-1 {
			pipeline.workers[idx+1].ChIn = pipeline.workers[idx].ChOut
		}
	}

	return pipeline, nil
}

func (p *pipeData) Send(val interface{}) {
	p.workers[0].ChIn <- val
}

func (p *pipeData) GetResChan() chan interface{} {
	return p.workers[len(p.workers)-1].ChOut
}

func (p *pipeData) Start() {
	for _, w := range p.workers {
		w.Start()
	}
}

func (p pipeData) Stop() {
	close(p.workers[0].ChIn)
	p.wg.Wait()
}
