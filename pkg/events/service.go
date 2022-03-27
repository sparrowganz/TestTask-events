package events

type Service interface {
	SendEvent(event []byte) error
}

type Repository interface {
	Save(data interface{}) error
}
