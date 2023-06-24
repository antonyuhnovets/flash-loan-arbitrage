package app

import (
	"context"
	"fmt"
	"log"

	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/contract"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/provider"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/repo"
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
		conf.Blockchain.Url,
		os.Getenv("ACCOUNT_PRIVATE_KEY"),
	)
	if err != nil {
		fmt.Println(err)
	}

	cAdress := common.HexToAddress(
		conf.Contract.Address,
	)

	c, err := cl.DialContract(cAdress)
	if err != nil {
		log.Fatal(err)
	}

	contract := contract.NewContract(
		cAdress, c, make([]entities.Token, 0),
	)
	provider := provider.NewTradeProvider(
		ctx, cl, entities.TradePair{},
	)
	repository := repo.UseFile("./storage_test/test.json")
	tc := tradecase.New(repository, provider, contract)

	// err = tc.Trade(ctx)
	log.Println(tc)
}

func icpConn(ctx context.Context, tc tradecase.TradeCase) {
	r := tc.GetRepo()
	apiRepo := *&rpc.API{
		Namespace:     "repo",
		Version:       "1.0",
		Service:       r,
		Public:        false,
		Authenticated: false,
	}
	api := make([]rpc.API, 0)
	api = append(api, apiRepo)

	ipcSrv, err := ethereum.NewEndpointIPC("test", api)
	if err != nil {
		log.Println(err)
	}
	go ipcSrv.Start()
	defer ipcSrv.Server.Stop()

	ipcCl, err := ipcSrv.AddClient(ctx)
	if err != nil {
		log.Println(err)
	}
	log.Println(ipcCl.Client.SupportedModules())
	cli := ipcSrv.DialInProc()
	log.Println(ipcCl.Client == cli.Client)
	log.Println(cli.Client.SupportedModules())
	newCl, err := rpc.DialContext(ctx, "test")

	i, err := repo.NewFile("./storage_test/test_rpc.json")

	newCl.RegisterName("repository", ipcSrv.Server)
	m, err := newCl.SupportedModules()
	log.Println(newCl.Call(i, m["rpc"]))
	log.Println(ipcSrv.Api)
}

func SetupTradeCase()

func IsDeployed(address string) bool {
	if address == "" {
		return false
	}

	return true
}

func Deploy(conf *config.Config) (string, error) {
	arg := fmt.Sprintf("network=%s", conf.Blockchain.Name)
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
	cmd := exec.Command("make", "build", conf.Contract.Name)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func Verify(conf *config.Config) {}
