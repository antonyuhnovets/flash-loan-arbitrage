package app

import (
	"context"
	"fmt"
	"log"
	"syscall"

	"os"
	"os/exec"
	"os/signal"

	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	v1 "github.com/antonyuhnovets/flash-loan-arbitrage/internal/delivery/rest/v1"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade/contract"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade/parser"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade/provider"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/trade/repo"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/httpserver"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

func Run(conf *config.Config) {

	ctx := context.Background()

	// ethereum client setup
	cl, err := ethereum.NewClient(
		conf.Blockchain.Url,
	)
	if err != nil {
		fmt.Println(err)
	}

	// contract connect
	// cAdress := common.HexToAddress(
	// 	conf.Blockchain.Contract.Address,
	// )
	cont, err := contract.New(conf.Blockchain.Contract.Address, cl)
	if err != nil {
		log.Fatal(err)
	}

	// Tradecase

	// contract instance
	ctr := contract.NewFlashArbContract(
		cont,
		make([]entities.TradePair, 0),
	)
	// provider create
	provider, err := provider.NewTradeProvider(
		ctx, conf.Blockchain.Url, os.Getenv("ACCOUNT_PRIVATE_KEY"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// repository setup
	var repository trade.Repository

	switch conf.Storage.Type {
	case "localfile":
		files := map[string]string{
			"pools": fmt.Sprintf(
				"%s/pools.json",
				conf.Storage.Localstorage.Path,
			),
			"tokens": fmt.Sprintf(
				"%s/tokens.json",
				conf.Storage.Localstorage.Path,
			),
		}
		repository, err = repo.NewStorage(files)
		if err != nil {
			log.Fatal(err)
		}
	case "database":
		if conf.Storage.Database.Driver == "postgres" {
			repository, err = repo.New(conf.Database)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// new tradecase
	tc := trade.New(
		repository,
		provider,
		ctr,
	)

	// store tokens
	// tokenPair := entities.TokenPair{
	// 	Token0: entities.Token{
	// 		ID:      1,
	// 		Name:    "WETH",
	// 		Address: "0xb4fbf271143f4fbf7b91a5ded31805e42b2208d6",
	// 		WeiVal:  1000000000000000000,
	// 	},
	// 	Token1: entities.Token{
	// 		ID:      2,
	// 		Name:    "LINK",
	// 		Address: "0x326C977E6efc84E512bB9C30f76E30c160eD06FB",
	// 		WeiVal:  1000000000000000000,
	// 	},
	// }

	// tc.Repo.AddToken(ctx, "tokens", tokenPair.Token0)
	// tc.Repo.AddToken(ctx, "tokens", tokenPair.Token1)

	// err = tc.SetTokens(ctx, "tokens")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pairList := make([]entities.TokenPair, 0)
	// pairList = append(pairList, tokenPair)

	// Parsecase

	// new parser with protocol
	p := parser.NewParser(entities.SwapProtocol{
		Name:       "Uniswap-V2",
		Factory:    "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f",
		SwapRouter: "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
	})
	// p["uniswap-v3"] = &parseUniV3

	// parsecase create
	pc := trade.NewParseCase(
		repository,
		p,
	)

	// logger
	l := logger.New(conf.Log.Level)

	// http server
	handler := gin.New()
	v1.NewRouter(handler, l, *tc, pc)
	httpServer := httpserver.New(
		handler,
		httpserver.Port(conf.HttpServer.Port),
	)

	// waiting signal
	interrupt := make(
		chan os.Signal,
		1,
	)
	signal.Notify(
		interrupt,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// run server
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf(
			"app - Run - httpServer.Notify: %w",
			err,
		))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf(
			"app - Run - httpServer.Shutdown: %w",
			err,
		))
	}

	// err = tc.Trade(ctx)
	// log.Println(tc)
}

func IsDeployed(address string) bool {

	return address != ""
}

func Deploy(conf *config.Config) (string, error) {
	arg := fmt.Sprintf(
		"network=%s",
		conf.Blockchain.Name,
	)
	arg1 := fmt.Sprintf(
		"contract=%s",
		conf.Contract.Address,
	)

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
	arg := fmt.Sprintf(
		"contract=%s",
		conf.Contract.Address,
	)

	cmd := exec.Command("make", "build", arg)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func Verify(conf *config.Config) {}
