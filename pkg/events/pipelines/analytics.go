package pipelines

import (
	"TestTask-events/config"
	"TestTask-events/pkg/app"
	"TestTask-events/pkg/workers"
	"encoding/json"
	"github.com/pkg/errors"
)

func CreateEventAnalyticsPipeline(
	core app.Core,
	config config.Workers,
) (workers.Pipeline, error) {

	//Create analytics pipe
	unmarshallJsonWorker := workers.New(core, config.Analytics, config.Buffer, unmarshallJsonWorkerFunc)
	setDataWorker := workers.New(core, config.Analytics, config.Buffer, setFieldsJsonWorkerFunc)

	analyticsPipe, err := workers.NewPipeline(unmarshallJsonWorker, setDataWorker)
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

	//Set more faster json unmarshaller (example easyjson)
	err := json.Unmarshal(val.([]byte), &data)
	if err != nil {
		return errors.Wrap(err, "failed unmarshallJsonWorkerFunc")
	}

	resChan <- data
	return nil
}

//Reflection is very bad idea
func setFieldsJsonWorkerFunc(val interface{}, resChan chan<- interface{}) error {
	data := val.(map[string]interface{})
	data["ip"] = "8.8.8.8"
	data["server_time"] = "2020-12-01 23:53:00"
	resChan <- data
	return nil
}
