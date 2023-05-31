package main

import (
	"os"

	"github.com/antonyuhnovets/flash-loan-arbitrage/scripts"
)

func main() {
	scripts.DeployAAVE(os.Getenv("RPC_URL"), os.Getenv("PRIVATE_KEY"), os.Getenv("ADDRESS_PROVIDER"))
}
