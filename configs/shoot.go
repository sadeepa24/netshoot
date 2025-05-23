package config

import "go.uber.org/zap"

type Config struct {
	Client Client     `json:"client"`
	Server Server     `json:"server"`
	Result Result     `json:"result"`
	Host   HostMgConf `json:"host"`
	Log    Logger     `json:"log"`
}

type Logger struct {
	Level    string `json:"level"`
	Output   []string `json:"paths"`
	Encoding string   `json:"encode"`
}

func (lg *Logger) ZapLevel() zap.AtomicLevel {
	switch lg.Level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.WarnLevel) 
	}
}