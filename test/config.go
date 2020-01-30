package test

import (
	"fmt"
	"github.com/kindritskyiMax/lets/config"
	"os"
)

func GetTestConfig() *config.Config {
	conf, err := config.Load("lets.yaml", "..")
	if err != nil {
		fmt.Printf("can not read test config: %s", err)
		os.Exit(1)
	}
	return conf
}
