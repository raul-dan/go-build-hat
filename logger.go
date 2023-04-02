package buildhat

import "go.uber.org/zap"

var logger, _ = zap.NewProduction()

func SetLogger(l *zap.Logger) {
	logger = l
}
