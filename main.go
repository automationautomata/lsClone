package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type config struct {
	Port string `json: port`
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

func checkInput(rootpath string, sort string) (bool, error) {
	if rootpath == "" {
		return false, errors.New("Укажите корневую директорию")

	}
	if _, err := os.Stat(rootpath); err != nil {
		fmt.Println(err.Error())
		return false, errors.New("Корневая директория не существует")
	}

	switch strings.ToLower(sort) {
	case "asc":
		return true, nil
	case "desc":
		return true, nil
	default:
		return false, errors.New("Указан неверный тип сортировки")
	}
}

func getArray(root string) ([]*lsCloneInfo, error) {
	eg := new(errgroup.Group)

	var table []*lsCloneInfo
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		lsInfo := NewlsCloneInfo(entry.Name(), entry.IsDir())
		if entry.IsDir() {
			eg.Go(func() error {
				return lsInfo.calcSize(filepath.Join(root, entry.Name()))
			})
		} else {
			err = lsInfo.IncreaseBy(info.Size())
			if err != nil {
				return nil, err
			}
		}
		table = append(table, lsInfo)
	}

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	return table, nil
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sortHeader := r.Form.Get("sort")
	if sortHeader == "" {
		sortHeader = "asc"
	}
	rootHeader := r.Form.Get("root")

	sortType, err := checkInput(rootHeader, sortHeader)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	table, err := getArray(rootHeader)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	sort.SliceStable(table, func(i, j int) bool {
		if sortType {
			return table[i].GetSize() < table[j].GetSize()
		} else {
			return table[i].GetSize() > table[j].GetSize()
		}
	})

	bytes, err := json.Marshal(table)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(bytes))
}

func main() {
	portFlag := flag.String("port", "", "Порт, на котором работает сервер")
	flag.Parse()
	fmt.Println("Запуск")
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var port string
	if *portFlag != "" {
		port = *portFlag
	} else {
		config, err := readConfig("config.json")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		port = config.Port
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.HandleFunc("/fs", handleQuery)

	fmt.Println("Сервер работает на порту", port)

	eg, egCtx := errgroup.WithContext(mainCtx)
	eg.Go(func() error {
		return srv.ListenAndServe()
	})
	eg.Go(func() error {
		<-egCtx.Done()
		return srv.Shutdown(context.Background())
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
