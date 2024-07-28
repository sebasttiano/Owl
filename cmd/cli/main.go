package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/sebasttiano/Owl/internal/cli"
	"github.com/sebasttiano/Owl/internal/config"
	"github.com/sebasttiano/Owl/internal/logger"
	"go.uber.org/zap"
)

var buildVersion = "N/A"
var buildDate = "N/A"

type templateInfoEntry struct {
	Version string
	Date    string
}

//go:embed client_info
var clientInfo string

func main() {
	if err := logger.Initialize("DEBUG"); err != nil {
		fmt.Println("logger initialization failed")
		return
	}

	cfg, err := config.NewClientConfig()
	if err != nil {
		logger.Log.Error("parsing config failed", zap.Error(err))
		return

	}

	tmpl, err := template.New("info").Parse(clientInfo)
	if err != nil {
		logger.Log.Error("failed to render banner: %v", zap.Error(err))
	}
	buf := new(bytes.Buffer)
	tmpl.Execute(buf, templateInfoEntry{buildVersion, buildDate})
	cfg.Info.Banner = buf.String()
	fmt.Println(cfg.Info.Banner)
	cliApp := cli.NewCLI(cfg)
	cliApp.Run()
}
