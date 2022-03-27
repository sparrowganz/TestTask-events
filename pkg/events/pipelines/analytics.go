package pipelines

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sparrowganz/TestTask-events/pkg/app"
	"github.com/sparrowganz/TestTask-events/pkg/events"
	"github.com/sparrowganz/TestTask-events/pkg/workers"
	"sync"
)

func CreateEventAnalyticsPipeline(core app.Core, repository events.Repository) (workers.Pipeline, error) {
	//todo set to config
	//Create analytics pipe
	wg := &sync.WaitGroup{}
	unmarshallJsonWorker := workers.New(wg, core, 1, 10, unmarshallJsonWorkerFunc)
	setDataWorker := workers.New(wg, core, 1, 10, setFieldsJsonWorkerFunc)
	saveWorker := workers.New(wg, core, 1, 10,
		func(val interface{}, resChan chan<- interface{}) error {
			return repository.Save(val)
		},
	)

	analyticsPipe, err := workers.NewPipeline(wg, unmarshallJsonWorker, setDataWorker, saveWorker)
	if err != nil {
		return nil, err
	}

	//Grace shutdown pipe
	analyticsPipe.Start()
	core.Group().Go(func() error {
		<-core.Context().Done()
		analyticsPipe.Stop()
		return nil
	})

	return analyticsPipe, nil
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