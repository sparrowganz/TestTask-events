package kafka

type kafka struct {
}

func New() *kafka {
	return &kafka{}
}

func (k *kafka) Send() error {
	return nil
}
