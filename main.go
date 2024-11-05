package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sync/errgroup"
)

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
		info, _ := entry.Info()

		var entryType EntryType
		if info.IsDir() {
			entryType = Folder
		} else {
			entryType = File
		}

		lsInfo := &lsCloneInfo{Name: entry.Name(), Type: entryType}

		if entry.IsDir() {
			eg.Go(func() error {
				return lsInfo.calcSize(filepath.Join(root, entry.Name()))
			})
		} else {
			err = lsInfo.IncreaseBy(info.Size())
			if err != nil {
				fmt.Println(err.Error())
				return nil, nil
			}
		}
		table = append(table, lsInfo)
	}

	err = eg.Wait()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return table, nil
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	sortHeader, rootHeader := "asc", ""
	if r.Form.Has("sort") {
		sortHeader = r.Form["sort"][0][1 : len(r.Form["sort"][0])-2]
	}
	if r.Form.Has("root") {
		rootHeader = r.Form["root"][0][1 : len(r.Form["root"][0])-2]
	}
	fmt.Printf("%#v\n", (rootHeader))

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
	// Get information about the file "main.go"
	fi, err := os.Stat("main.go")
	if err != nil {
		// Handle the error
		fmt.Println(err)
		return
	}
	// Print the file's size
	fmt.Println("File size:", fi.Size(), fi.Name())

	http.HandleFunc("/fs", handleQuery)     // устанавливаем обработчик
	err = http.ListenAndServe(":8086", nil) // устанавливаем порт, который будем слушать
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
