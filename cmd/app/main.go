package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/telepuz/alerton/internal/alert"
	"github.com/telepuz/alerton/internal/app"
	config "github.com/telepuz/alerton/internal/config"
	"github.com/telepuz/alerton/internal/logger"
	"github.com/telepuz/alerton/internal/messenger/telegram"
	"github.com/telepuz/alerton/internal/storage/memory"
)

var (
	scriptsDir  *string
	configFile  *string
	versionFlag *bool
	version     = "devel"
	gitRevision = "devel"
	buildDate   = "devel"
)

func init() {
	scriptsDir = flag.String("scripts_dir", "/etc/alerton/scripts", "Scripts directory")
	configFile = flag.String("config_file", "/etc/alerton/alerton.yml", "Config filename")
	versionFlag = flag.Bool("version", false, "Print version")
}

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf(
			"Name: %s\nVersion: %s\ngitRevision: %s\nbuildDate: %s\n",
			"AlertON",
			version,
			gitRevision,
			buildDate)
		os.Exit(0)
	}

	cfg, err := config.NewConfig(*configFile)
	if err != nil {
		slog.Error(
			fmt.Sprintf("main(): %s", err))
		os.Exit(1)
	}
	slog.Debug("main(): Read config file")

	err = logger.ConfigureSlog(&cfg.Logger)
	if err != nil {
		slog.Error(
			fmt.Sprintf("main(): %s", err))
		os.Exit(1)
	}
	slog.Debug("main(): Configured slog")
	slog.Debug(fmt.Sprintf("main(): Read configs: %+v", cfg))

	alerts, err := alert.NewAlerts(&cfg.Alerts, *scriptsDir)
	if err != nil {
		slog.Error(
			fmt.Sprintf("main(): %s", err))
		os.Exit(1)
	}
	slog.Debug("main(): Created alerts")

	msg, err := telegram.NewTelegram(&cfg.Messenger)
	if err != nil {
		slog.Error(
			fmt.Sprintf("main(): %s", err))
		os.Exit(1)
	}
	slog.Debug("main(): Created new messanger")

	storage := memory.NewMemoryStorage(&cfg.Storage)
	slog.Debug("main(): Created storage")

	appContext := app.AppContext{
		Config:    cfg,
		Messenger: msg,
		Alerts:    alerts,
		Storage:   storage,
	}

	slog.Info("main(): Starting app...")
	app.Run(&appContext)
	slog.Info("main(): Exit...")
}
