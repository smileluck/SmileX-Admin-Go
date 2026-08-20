// Package logger 全局 zap 日志
package logger

import (
	"go.uber.org/zap"
)

var L *zap.Logger

func Init(mode string) error {
	var (
		l   *zap.Logger
		err error
	)
	if mode == "release" {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		return err
	}
	L = l
	return nil
}

func Debug(msg string, fields ...zap.Field) { L.Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field)  { L.Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field)  { L.Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { L.Error(msg, fields...) }
