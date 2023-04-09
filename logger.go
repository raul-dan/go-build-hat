package buildhat

import (
	"buildhat/logger"
	"go.uber.org/zap"
)

func SetLogger(l *zap.Logger) {
	logger.SetInstance(l)
}
