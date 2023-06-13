package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
)

func App(conf *config.Config) {
	out, err := Deploy(conf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deployed")
	log.Println(out)
}

func Deploy(conf *config.Config) (string, error) {
	arg := fmt.Sprintf("deploy-%s", conf.Blockchain.Network)
	cmd := exec.Command("make", arg)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return "Failed with ", err
	}

	return "Success ", nil
}
