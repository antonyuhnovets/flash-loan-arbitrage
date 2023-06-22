package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
)

func App(conf *config.Config) {

	// if !IsDeployed(conf.Contract.Address) {
	// 	address, err := Deploy(conf)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("Deployed to: %s", address)
	// 	conf.Blockchain.Contract.Address = address
	// }

	// ctx := context.TODO()
	// cmd1 := cli.NewMakeCMD(
	// 	"build",
	// 	&ctx,
	// 	fmt.Sprintf("contract=%s", conf.Blockchain.Contract.Name),
	// )
	// cmd2 := cli.NewMakeCMD(
	// 	"deploy",
	// 	&ctx,
	// 	fmt.Sprintf("network=%s", conf.NetworkChain.Name),
	// )
	// cmd3 := cli.NewMakeCMD(
	// 	"delete",
	// 	&ctx,
	// 	fmt.Sprintf("contract=%s", conf.Blockchain.Contract.Name),
	// )
	// cmds := make(map[string]*cli.Command)
	// cmds["build"] = cmd1
	// cmds["deploy"] = cmd2
	// cmds["delete"] = cmd3

	// cli := cli.NewCLI(cmds)
	// cli.EnvSet("CONTRACT_ADDRESS", conf.Contract.Address)
	// cli.EnvSet("CONTRACT_NAME", conf.Blockchain.Contract.Name)

	// go cli.Run()
	// scan := cli.GetScanner()
	// for {
	// 	scan.Scan()
	// 	input := scan.Text()
	// 	cli.Cmd <- string(input)
	// }
	ctx := context.Background()
	cl, err := ethereum.NewClient(
		conf.Blockchain.NetworkChain.Url,
		os.Getenv("ACCOUNT_PRIVATE_KEY"),
	)
	if err != nil {
		fmt.Println(err)
	}

	cl.DialContract(ctx, conf.Blockchain.Contract.Address)
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

func Build(conf *config.Config) error {
	cmd := exec.Command("make", "build", conf.Blockchain.Contract.Name)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Verify(conf *config.Config) {}
