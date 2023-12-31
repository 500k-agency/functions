package main

import (
	"flag"
	"log"
	"os"

	// Blank-import the function package so the init() runs
	_ "github.com/500k-agency/function"
	"github.com/500k-agency/function/config"
	"github.com/500k-agency/function/lib/connect"
	"github.com/500k-agency/function/product"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

var (
	flags    = flag.NewFlagSet("main", flag.ExitOnError)
	confFile = flags.String("config", "", "path to config file")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("invalid flags: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside of your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
		conf, err := config.NewFromFile(*confFile, os.Getenv("CONFIG"))
		if err != nil {
			log.Fatalf("main.NewFromConfig: %v\n", err)
		}
		connect.Configure(conf.Connect)
		product.Setup(conf.Products)
	}

	log.Printf("server running on %s:%s", hostname, port)
	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}
