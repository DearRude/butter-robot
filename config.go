package main

import (
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	AppId    int
	AppHash  string
	BotToken string
	Logger   zap.Logger
}

func GenConfig() Config {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(zap.InfoLevel)
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Unable to make zap logger. Error: %s", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Unstable to sync logger")
		}
	}()

	logger.Info("Reading configurations")

	fs := flag.NewFlagSet("mastodon_exporter", flag.ContinueOnError)
	var (
		appId    = fs.Int("appId", 0, "APP ID in Telegram API")
		appHash  = fs.String("appHash", "", "APP Hash in Telegram API")
		botToken = fs.String("botToken", "", "Bot token given by botfather")
	)

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fs.String("config", "", "config file")
	} else {
		fs.String("config", ".env", "config file")
	}

	err = ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.EnvParser),
	)
	if err != nil {
		logger.Fatal("Unable to parse args")
	}

	return Config{
		AppId:    *appId,
		AppHash:  *appHash,
		BotToken: *botToken,
		Logger:   *logger,
	}
}
