package message

type Broker interface {
	Send() error
}
