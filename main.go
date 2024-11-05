package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
)

// calcSize - вычисление размерности папки.
func calcSize(lsInfo *lsCloneInfo, path string, eg *errgroup.Group) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			eg.Go(func() error {
				return calcSize(lsInfo, filepath.Join(path, entry.Name()), eg)
			})
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}

			lsInfo.IncreaseBy(info.Size())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// showFileInfo - выводит на экран информацию,
// sortType - тип сортировки: true - по возрастанию, false - по убыванию
func showFileInfo(table []*lsCloneInfo, sortType bool) {
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

func checkInput(rootpath string, sort string) (bool, error) {
	if rootpath == "" {
		return false, errors.New("Укажите корневую директорию")

	}
	if _, err := os.Stat(rootpath); err != nil {
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

func main() {
	rootFlag := flag.String("root", "", "Корневая директория")
	sortFlag := flag.String("sort", "ASC", "Тип сортировки: asc/desc")
	flag.Parse()

	sortType, err := checkInput(*rootFlag, *sortFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	eg := new(errgroup.Group)

	entries, err := os.ReadDir(*rootFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	var table []*lsCloneInfo

	for _, entry := range entries {
		info, _ := entry.Info()
		lsInfo := &lsCloneInfo{Name: entry.Name(), IsDir: info.IsDir(), size: 0}

		if entry.IsDir() {
			eg.Go(func() error {
				return calcSize(lsInfo, filepath.Join(*rootFlag, entry.Name()), eg)
			})
		} else {
			err = lsInfo.IncreaseBy(info.Size())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		table = append(table, lsInfo)
	}

	err = eg.Wait()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	showFileInfo(table, sortType)
}
