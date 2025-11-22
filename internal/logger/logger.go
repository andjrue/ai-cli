// Package logger provides unstructured logs - thought about slog but I don't need that
package logger

import (
	"log"
	"os"
)

var Log *log.Logger

func Init(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	
	Log = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}