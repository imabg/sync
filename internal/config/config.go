package config

import "go.uber.org/zap"

type Application struct {
	ErrorLog *zap.SugaredLogger
	InfoLog *zap.SugaredLogger
}