package main

import (
	"avito_project/cmd/app"
	"flag"
)

func main() {
	cfg := flag.String("config", "/config-dev.yml", "config file path")
	err := app.Run(*cfg)
	if err != nil {
		panic(err)
	}
}
