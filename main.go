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

func calcSize(lsinfo *lsCloneInfo, path string, wg *sync.WaitGroup, errChan chan error) error {
	defer wg.Done()

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var inner_wg sync.WaitGroup
	for _, entry := range entries {
		if entry.IsDir() {
			inner_wg.Add(1)
			go func() {
				err := calcSize(lsinfo, filepath.Join(path, entry.Name()), &inner_wg, errChan)
				errChan <- err
			}()
		} else {
			info, _ := entry.Info()
			lsinfo.Extend(info.Size())
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
		break
	case "ASC":
		sortType = true
		break
	case "DESC":
		sortType = false
		break
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
		lsInfo := &lsCloneInfo{Name: entry.Name(), IsDir: info.IsDir(), Size: 0}

		if entry.IsDir() {
			wg.Add(1)
			go func() {
				err = calcSize(lsInfo, filepath.Join(*rootFlag, entry.Name()), &wg, errChan)
				if err != nil {
					fmt.Println(err.Error())
				}
			}()
		} else {
			lsInfo.Extend(info.Size())
		}
		table = append(table, lsInfo)

	}

	go func() {
		if err := <-errChan; err != nil {
			fmt.Printf("Произошла ошибка: %v\n", err)
			cancel() // Отмена всех горутин
		}
	}()

	wg.Wait()

	sort.SliceStable(table, func(i, j int) bool {
		if sortType {
			return table[i].Size < table[j].Size
		} else {
			return table[i].Size > table[j].Size
		}
	})

	var entryType string
	for _, entry := range table {
		if entry.IsDir {
			entryType = "Folder"
		} else {
			entryType = "File"
		}
		fmt.Println(entryType, entry.Name, entry.convertSize(2))
	}
}
