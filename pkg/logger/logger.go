package logger

import "go.uber.org/zap"

type Logger struct {
	ErrorLog *zap.SugaredLogger
	InfoLog  *zap.SugaredLogger
	PanicLog *zap.SugaredLogger
}

var Log *Logger

func Setup() error {
	logger, err := zap.NewProduction()
	sugar := logger.Sugar()
	Log = &Logger{
		ErrorLog: sugar.Named("Error"),
		InfoLog:  sugar.Named("Info"),
		PanicLog: sugar.Named("Panic"),
	}
	defer logger.Sync()
	return err
}
