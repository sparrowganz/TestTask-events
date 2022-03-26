package message

type Broker interface {
	Send(value interface{}) error
}
