package main

import (
	"strconv"
	"fmt"
	"wdh-auth/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"
	"flag"
	"golang.org/x/net/context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

func main() {
	logger.Out = os.Stdout;

	serverPort, err := loadPortNumber()
	if err != nil {
		panic(err)
	}

	errChan := make(chan error, 2)

	var (
		httpAddr = flag.String("http.addr", ":"+strconv.Itoa(serverPort), "HTTP listen address")
	)

	var server *http.Server
	{
		http.HandleFunc("/", handlerCheckBot)
		server = &http.Server{
			Addr: *httpAddr,
		}

		go func() {
			logger.WithFields(logrus.Fields{
				"address": *httpAddr,
			}).Info("Http service listening")
			errChan <- server.ListenAndServe()
		}()
	}

	go func() {
		signals := make(chan os.Signal)
		signal.Notify(signals, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-signals)
	}()

	logger.Error("terminated", <-errChan)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	server.Shutdown(ctx)
}

func loadPortNumber() (int, error) {
	_, err := os.Stat(".env")

	if err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			return 0, fmt.Errorf("error loading .env file, %v", err)
		}
	}

	serverPort, err := strconv.Atoi(utils.EnvString("PORT", "80"))
	if err != nil {
		return 0, fmt.Errorf("can't read env params PORT")
	}

	return serverPort, nil
}

func handlerCheckBot(w http.ResponseWriter, r *http.Request) {
	detect := go_botdetect.NewBotDetect(r, nil)

	if detect.IsBot() {
		http.Error(w, "is bot", http.StatusUnauthorized)
		return
	}
}
