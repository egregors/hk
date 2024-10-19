package log

import (
	"io"
	"log"
	"os"
)

var (
	Info = &Logger{log.New(os.Stdout, "INFO", log.LstdFlags|log.Lshortfile)}
	Erro = &Logger{log.New(os.Stderr, "ERRO", log.LstdFlags|log.Lshortfile)}
)

type Logger struct {
	*log.Logger
}

func (l *Logger) On() {
	l.SetOutput(os.Stdout)
}

func (l *Logger) Off() {
	l.SetOutput(io.Discard)
}
