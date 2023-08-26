package main

import (
	"avito_project/cmd/app"
	"flag"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "config", "/config-dev.yml", "config file path")
	flag.Parse()
	err := app.Run(cfg)
	if err != nil {
		panic(err)
	}
}
