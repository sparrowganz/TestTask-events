package main

import (
	pkgMainApp "TestTask-events/pkg/app/web"
	"context"
	"flag"
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"

	"TestTask-events/config"
	"TestTask-events/pkg/app"
	"github.com/pkg/errors"
)

func main() {

	// Set path to configs file
	pathToConfigDir := os.Getenv(config.Env)
	if pathToConfigDir == "" {
		//Read params from OS
		flag.StringVar(&pathToConfigDir, "config", "", "Path to config file")
		flag.Parse()
	}

	// Read configs from environment or config file
	c, err := config.Parse(pathToConfigDir)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed read config"))
	}

	core := app.NewCore(context.Background(), c.App)
	core.Logger().Println("starting service...")

	//Graceful Shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		core.Logger().Println("stopping services...")
		core.Stop()
	}()

	//Create APP
	a, err := pkgMainApp.NewMain(c, core)
	if err != nil {
		core.Logger().Fatal(err.Error())
	}
	a.Start()

	//Wait stop APP
	if err := core.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
