package events

type Service interface {
	SendEvent(event []byte) error
}

type Repository interface {
	Save(event interface{}) error
}
