package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
)

func App(conf *config.Config) {

	if !IsDeployed(conf.Contract.Address) {
		address, err := Deploy(conf)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Deployed to: %s", address)
		conf.Blockchain.Contract.Address = address
	}

	go Cycle()
}

func IsDeployed(address string) bool {
	if address == "" {
		return false
	}

	return true
}

func Deploy(conf *config.Config) (string, error) {
	arg := fmt.Sprintf("network=%s", conf.NetworkChain.Name)
	arg1 := fmt.Sprintf("contract=%s", conf.Contract.Address)

	cmd := exec.Command("make", "deploy", arg, arg1)
	cmd.Stderr = os.Stderr

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	address := string(out)[len(string(out))-43:]
	os.Setenv("CONTRACT_ADDRESS", address)

	return address, nil
}

func Verify(conf *config.Config) {}

func Cycle() {
	for {

	}
}
