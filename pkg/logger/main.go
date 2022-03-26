package logger

import (
	"fmt"
	"log"
	"os"
)

//todo implement me
func New(appName string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf(" [%s] ", appName), 0)
}
