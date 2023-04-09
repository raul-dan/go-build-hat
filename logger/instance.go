package logger

import "go.uber.org/zap"

var Instance, _ = zap.NewProduction()

func SetInstance(l *zap.Logger) {
	Instance = l
}
