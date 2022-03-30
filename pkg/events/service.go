package events

type Service interface {
	SendEvent(event *Event) error
}

type Repository interface {
	Collect(event *Event)
}

type Event struct {
	ClientTime string `json:"client_time"`
	DeviceID   string `json:"device_id"`
	DeviceOS   string `json:"device_os"`
	Session    string `json:"session"`
	Sequence   int    `json:"sequence"`
	Event      string `json:"event"`
	ParamInt   int    `json:"param_int"`
	ParamStr   string `json:"param_str"`
	IP         string `json:"ip"`
	ServerTime string `json:"server_time"`
}
