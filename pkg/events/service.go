package events

type Service interface {
	SendEvent(event []byte) error
}
