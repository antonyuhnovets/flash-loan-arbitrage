package main

import (
	"log"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/app"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(
		"network: %s, \naccount: %s \ncontract: %s - %s\n",
		conf.Blockchain.Name,
		conf.Account.Address,
		conf.Contract.Name,
		conf.Contract.Address,
	)

	app.Run(conf)
}
