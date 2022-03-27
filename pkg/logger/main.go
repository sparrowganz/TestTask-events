package logger

import (
	"fmt"
	"log"
	"os"
)

//TODO set other logger
func New(appName string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf(" [%s] ", appName), 0)
}
