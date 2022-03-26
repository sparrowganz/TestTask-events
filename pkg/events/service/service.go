package service

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sparrowganz/TestTask-events/pkg/app"
	"github.com/sparrowganz/TestTask-events/pkg/workers"
	"log"
	"sync"
)

type Service struct {
	pipe workers.Pipeline
}

func New(core app.Core) (*Service, error) {

	//todo set to config
	wg := &sync.WaitGroup{}
	unmarshallJsonWorker := workers.New(wg, core, 1, 10, unmarshallJsonWorkerFunc)
	setDataWorker := workers.New(wg, core, 1, 10, setFieldsJsonWorkerFunc)
	resWorker := workers.New(wg, core, 1, 10, resWorker)

	pipe, err := workers.NewPipeline(wg, unmarshallJsonWorker, setDataWorker, resWorker)
	if err != nil {
		return nil, err
	}

	pipe.Start()
	core.Group().Go(func() error {
		<-core.Context().Done()
		pipe.Stop()
		return nil
	})

	return &Service{
		pipe: pipe,
	}, nil
}

func (s Service) SendEvent(event []byte) error {
	s.pipe.Send(event)
	return nil
}

func unmarshallJsonWorkerFunc(val interface{}, resChan chan<- interface{}) error {

	var data = make(map[string]interface{}, 10)

	//todo Faster
	err := json.Unmarshal(val.([]byte), &data)
	if err != nil {
		return errors.Wrap(err, "failed unmarshallJsonWorkerFunc")
	}

	resChan <- data
	return nil
}

//todo решить проблему с типами
func setFieldsJsonWorkerFunc(val interface{}, resChan chan<- interface{}) error {
	data := val.(map[string]interface{})
	data["ip"] = "8.8.8.8"
	data["server_time"] = "2020-12-01 23:53:00"
	resChan <- data
	return nil
}

func resWorker(val interface{}, resChan chan<- interface{}) error {
	log.Println(val)
	return nil
}
