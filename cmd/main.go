package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kingpin"
	"github.com/sadeepa24/netshoot"
	"github.com/sadeepa24/netshoot/cmd/tools"
	config "github.com/sadeepa24/netshoot/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	versionNum   = ""
)

var (
	app = kingpin.New("netshoot", "A versatile tool for discovering hosts and identifying tunnelable ones with many customizable options.")
	run        = app.Command("run", "run netshoot")
	configDir = run.Flag("config", "config path").
			Short('c').
			Default("config.json").
			String()
	version     = app.Command("version", "netshoot version")
	gen = app.Command("generate", "generate payload file")
	pfiledir = gen.Flag("file", "payload info file").Short('i').Required().String()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
		case run.FullCommand():
			netshoot_run()
		case version.FullCommand():
			fmt.Println(versionNum)
		case gen.FullCommand():
			err := tools.GenPayloadFile(*pfiledir)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("payload file created.")
	}

}

func netshoot_run() {

	configb, err := os.ReadFile(*configDir)
	if err != nil {
		log.Fatal("config file open err: " + err.Error())
	}
	var config config.Config
	err = json.Unmarshal(configb, &config)
	if err != nil {
		log.Fatal("config Unmarshal Err: " + err.Error())
	}

	logconfig := zap.Config{
		Level: config.Log.ZapLevel(),
		Development: false,
		DisableCaller: true,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: config.Log.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       config.Log.ZapLevel().Level().CapitalString(),
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			// StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: config.Log.Output,
		ErrorOutputPaths: config.Log.Output,
	}
	logger, err := logconfig.Build()
	if err != nil {
		log.Fatal("logger build fail: " + err.Error())
	}
	ctx, cancel := context.WithCancel(context.Background())
	shoot, err := netshoot.New(ctx, logger, config)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = shoot.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	shoot.Close()
	cancel()
}