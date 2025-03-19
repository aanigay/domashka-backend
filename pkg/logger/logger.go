package logger

import "context"

// TODO: Прикрутить нормальный логгер

type SomeLogger struct {
}

func New() *SomeLogger {
	return &SomeLogger{}
}

func (l *SomeLogger) Error(_ context.Context, _ ...interface{}) {

}

func (l *SomeLogger) Info(_ context.Context, _ ...interface{}) {

}
