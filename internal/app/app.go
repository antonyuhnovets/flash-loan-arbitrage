package app

import (
	"context"
	"fmt"
	"log"
	"syscall"

	"os"
	"os/exec"
	"os/signal"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/antonyuhnovets/flash-loan-arbitrage/config"
	v1 "github.com/antonyuhnovets/flash-loan-arbitrage/internal/delivery/rest/v1"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/entities"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/contract"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/parser"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/provider"
	"github.com/antonyuhnovets/flash-loan-arbitrage/internal/tradecase/repo"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/ethereum"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/httpserver"
	"github.com/antonyuhnovets/flash-loan-arbitrage/pkg/logger"
)

func Run(conf *config.Config) {

	ctx := context.Background()

	// ethereum client
	cl, err := ethereum.NewClient(
		conf.Blockchain.Url,
		os.Getenv("ACCOUNT_PRIVATE_KEY"),
	)
	if err != nil {
		fmt.Println(err)
	}

	// contract
	cAdress := common.HexToAddress(
		conf.Contract.Address,
	)
	ctr, err := cl.DialContract(cAdress)
	if err != nil {
		log.Fatal(err)
	}

	// tradecase
	contract := contract.NewContract(
		cAdress, ctr,
		make([]entities.TradePair, 0),
	)
	provider := provider.NewTradeProvider(
		ctx, cl,
	)
	repository := repo.NewStorage()
	repository.UseFile(
		"pools",
		"./storage_test/pools_test.json",
	)
	repository.UseFile(
		"tokens",
		"./storage_test/tokens_test.json",
	)
	tc := tradecase.New(
		repository,
		provider,
		contract,
	)

	tokenPair := entities.TokenPair{
		Token0: entities.Token{
			Name:    "WETH",
			Address: "0xb4fbf271143f4fbf7b91a5ded31805e42b2208d6",
			WeiVal:  1000000000000000000,
		},
		Token1: entities.Token{
			Name:    "LINK",
			Address: "0x326C977E6efc84E512bB9C30f76E30c160eD06FB",
			WeiVal:  1000000000000000000,
		},
	}

	tc.Repo.AddToken(ctx, "tokens", tokenPair.Token0)
	tc.Repo.AddToken(ctx, "tokens", tokenPair.Token1)

	err = tc.SetTokens(ctx, "tokens")

	pairList := make([]entities.TokenPair, 0)
	pairList = append(pairList, tokenPair)

	// parseUni := parser.NewUniV2(entities.SwapProtocol{
	// 	Name:       "Uniswap-V2",
	// 	Factory:    "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f",
	// 	SwapRouter: "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45",
	// },
	// 	pairList,
	// )
	// tc.SetParser(&parseUni)
	// tc.ParseWrite(ctx, "pools", pairList)

	// parseSushi := parser.NewSushiV2(entities.SwapProtocol{
	// 	Name:       "Sushiswap-V2",
	// 	Factory:    "0xc35DADB65012eC5796536bD9864eD8773aBc74C4",
	// 	SwapRouter: "0x1b02dA8Cb0d097eB8D57A175b88c7D8b47997506",
	// },
	// 	pairList,
	// )
	// tc.SetParser(&parseSushi)
	// tc.ParseWrite(ctx, "pools", pairList)

	parseUniV3 := parser.NewUniV3(
		pairList,
		entities.SwapProtocol{
			Name:       "Uniswap-V3",
			Factory:    "0x1F98431c8aD98523631AE4a59f267346ea31F984",
			SwapRouter: "0xE592427A0AEce92De3Edee1F18E0157C05861564",
		},
	)

	p := make(map[string]tradecase.Parser)
	p["uniswap-v3"] = &parseUniV3
	pc := tradecase.NewParseCase(
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

	// run
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
	if address == "" {
		return false
	}

	return true
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
