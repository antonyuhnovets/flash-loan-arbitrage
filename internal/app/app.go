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
	c, err := cl.DialContract(cAdress)
	if err != nil {
		log.Fatal(err)
	}

	// tradecase
	contract := contract.NewContract(
		cAdress, c, make([]entities.Token, 0),
	)
	provider := provider.NewTradeProvider(
		ctx, cl, entities.TradePair{},
	)
	repository := repo.UseFile("./storage_test/test.json")
	tc := tradecase.New(repository, provider, contract)

	// logger
	l := logger.New(conf.Log.Level)

	// http server
	handler := gin.New()
	v1.NewRouter(handler, l, *tc)
	httpServer := httpserver.New(
		handler,
		httpserver.Port(conf.HttpServer.Port),
	)

	// waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// run
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
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
	arg := fmt.Sprintf("contract=%s", conf.Contract.Address)

	cmd := exec.Command("make", "build", arg)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func Verify(conf *config.Config) {}
