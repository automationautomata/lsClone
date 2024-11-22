package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

func checkPath(path string) error {
	if path == "" {
		return errors.New("Укажите директорию")
	}
	if _, err := os.Stat(path); err != nil {
		fmt.Println(err.Error())
		return errors.New("Директория не существует")
	}
	return nil
}

func getFolderData(root string, sortType string) ([]*EntryInfo, error) {
	var isAscSort bool
	switch strings.ToLower(sortType) {
	case "asc":
		isAscSort = true
	case "":
		isAscSort = true
	case "desc":
		isAscSort = false
	default:
		error_text := "Указан неверный тип сортировки"
		log.Println(error_text)
		return nil, errors.New(error_text)
	}

	err := checkPath(root)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	table, err := createEntriesTable(root)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	for i := range table {
		table[i].convertSize(2)
	}

	sort.SliceStable(table, func(i, j int) bool {
		if isAscSort {
			return table[i].GetSize() < table[j].GetSize()
		} else {
			return table[i].GetSize() > table[j].GetSize()
		}
	})

	return table, nil
}

type ViewData struct {
	Separator string
	Folders   []string
}

func createQueryHandler(staticsDir string, startRoot string, statisticsServerAddr string, htmlPath string) *http.ServeMux {
	mux := http.NewServeMux()
	tpl := template.Must(template.ParseFiles(htmlPath))

	fs := http.FileServer(http.Dir(staticsDir))
	staticsPrefix := fmt.Sprint("/", filepath.Base(staticsDir), "/")
	mux.Handle(staticsPrefix, http.StripPrefix(staticsPrefix, fs))

	var sep = string(filepath.Separator)
	folders := strings.Split(startRoot, sep)
	if folders[0] == "" {
		folders = folders[1:]
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Execute(w, &ViewData{sep, folders})
		if err != nil {
			log.Println(err.Error())
		}
	})
	mux.HandleFunc("/fs", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("Content-Type", "application/json")

		r.ParseForm()
		rootHeader, sortHeader := "", ""
		if r.Form.Has("root") {
			rootHeader = r.Form.Get("root")
		}
		if r.Form.Has("sort") {
			sortHeader = r.Form.Get("sort")
		}

		table, err := getFolderData(rootHeader, sortHeader)
		if err != nil {
			fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", err.Error(), "\" }"))
			return
		} else {
			bytes, err := json.Marshal(table)
			if err != nil {
				log.Println(err.Error())
				fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", err.Error(), "\" }"))
				return
			}
			fmt.Fprintln(w, string(bytes))
		}

		timeDelta := time.Since(start)
		go func() {
			sumSize := int64(0)
			for _, entry := range table {
				sumSize += entry.GetSize()
			}
			data := fmt.Sprint("{\"path\": \"", rootHeader, "\", \"size\": \"",
				sumSize, "\", \"time\": \"", timeDelta.Nanoseconds(), "\" }")
			resp, err := http.Post(statisticsServerAddr, "application/json", bytes.NewBuffer([]byte(data)))
			if err != nil {
				log.Println(err.Error())
				return
			}

			body, _ := io.ReadAll(resp.Body)
			log.Println(fmt.Sprint("Stat response Status:", resp.Status,
				"Headers:", resp.Header, "Body:", string(body)))
		}()
	})

	return mux
}
