package app

type Saver interface {
}

type saverData struct {
}

func NewSaver() Saver {
	return &saverData{}
}
