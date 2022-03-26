package kafka

type Config struct {
	Address string `yaml:"address"`
	Topic   string `yaml:"events"`
	Group   string `yaml:"group"`
}

type kafka struct {
}

func New() *kafka {
	return &kafka{}
}

func (k *kafka) Send() error {
	return nil
}
