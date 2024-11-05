package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// calcSize - вычисление размерности папки.
func calcSize(lsinfo *lsCloneInfo, path string, wg *sync.WaitGroup, errChan chan error) error {
	defer wg.Done()

	entries, err := os.ReadDir(path)
	if err != nil {
		errChan <- err
		return err
	}

	var inner_wg sync.WaitGroup
	for _, entry := range entries {
		if entry.IsDir() {
			inner_wg.Add(1)
			go calcSize(lsinfo, filepath.Join(path, entry.Name()), &inner_wg, errChan)

		} else {
			info, err := entry.Info()
			if err != nil {
				errChan <- err
			}

			lsinfo.IncreaseBy(info.Size())
			if err != nil {
				errChan <- err
			}
		}
	}

	inner_wg.Wait()
	return nil
}

func main() {
	rootFlag := flag.String("root", "", "Корневая директория")
	sortFlag := flag.String("sort", "", "Тип сортировки: asc/desc")
	flag.Parse()

	if _, err := os.Stat(*rootFlag); err != nil {
		fmt.Println("Корневая директория не существует")
		return
	}

	var sortType bool
	switch strings.ToUpper(*sortFlag) {
	case "":
		sortType = true
	case "ASC":
		sortType = true
	case "DESC":
		sortType = false
	default:
		fmt.Println("Указан неверный тип сортировки")
		return
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	entries, err := os.ReadDir(*rootFlag)
	if err != nil {
		log.Fatal(err)
	}

	var table []*lsCloneInfo
	var wg sync.WaitGroup
	errChan := make(chan error)

	for _, entry := range entries {
		info, _ := entry.Info()
		lsInfo := &lsCloneInfo{Name: entry.Name(), IsDir: info.IsDir(), size: 0}

		if entry.IsDir() {
			wg.Add(1)
			go calcSize(lsInfo, filepath.Join(*rootFlag, entry.Name()), &wg, errChan)
		} else {
			err = lsInfo.IncreaseBy(info.Size())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		table = append(table, lsInfo)

	}

	go func() {
		if err := <-errChan; err != nil {
			fmt.Printf("Произошла ошибка: %v\n", err)
			cancel()
		}
	}()

	wg.Wait()

	sort.SliceStable(table, func(i, j int) bool {
		if sortType {
			return table[i].GetSize() < table[j].GetSize()
		} else {
			return table[i].GetSize() > table[j].GetSize()
		}
	})

	for _, entry := range table {
		if entry.IsDir {
			fmt.Print("Folder ")
		} else {
			fmt.Print("File ")
		}
		fmt.Println(entry.Name, entry.convertSize(2))

	}
}
