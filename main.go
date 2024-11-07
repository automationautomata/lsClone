package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type config struct {
	Port     string `json:"port"`
	LogsPath string `json:"logspath"`
}

func readConfig(path string) (*config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("Файл конфигурации не найден")
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &config{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, errors.New("Файл конфигурации не найден")
	}
	return config, nil
}

func main() {
	portFlag := flag.String("port", "", "Порт, на котором работает сервер")
	flag.Parse()
	fmt.Println("Запуск")
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := readConfig("config.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	port := config.Port

	file, err := os.OpenFile(config.LogsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		fmt.Println("Не удалось открыть файл логов, используется стандартный stderr")
	}
	log.Println("Запуск")

	if *portFlag != "" {
		port = *portFlag
	}

	mux := createQueryHandler("ui/static", "ui/html/index.html")
	srv := &http.Server{
		Addr:    fmt.Sprint(":", port),
		Handler: mux,
	}
	fmt.Println("Сервер работает на порту", port)
	log.Println("Сервер работает на порту", port)

	eg, egCtx := errgroup.WithContext(mainCtx)
	eg.Go(func() error {
		return srv.ListenAndServe()
	})
	eg.Go(func() error {
		<-egCtx.Done()
		return srv.Shutdown(context.Background())
	})

	if err := eg.Wait(); err != nil {
		fmt.Println("exit reason:", err)
		log.Fatalln("exit reason:", err)
	}
}
