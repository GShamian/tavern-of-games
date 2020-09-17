package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"

	"github.com/GShamian/tavern-of-games/internal/app/apiserver"
)

var (
	configPath string
)

// init func. Initialising CLI flags. If we execute program
// without flags it automatically sets default values.
func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

// main func. Starts server with config settings
func main() {
	// Parsing flags
	flag.Parse()
	// Creating config entity for server
	config := apiserver.NewConfig()
	// Decoding toml config file
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	// Starting server with our config
	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
