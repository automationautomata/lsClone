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
	"text/template"
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

func getEntriesTableHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	var sortType bool
	if r.Form.Has("sort") {
		switch strings.ToLower(r.Form.Get("sort")) {
		case "asc":
			sortType = true
		case "":
			sortType = true
		case "desc":
			sortType = false
		default:
			error_text := "Указан неверный тип сортировки"
			fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", error_text, "\" }"))
			log.Println(w, error_text)
			return
		}
	} else {
		sortType = true
	}

	rootHeader := r.Form.Get("root")

	err := checkPath(rootHeader)
	if err != nil {
		fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", err.Error(), "\" }"))
		log.Println(w, err.Error())
		return
	}

	table, err := getEntriesTable(rootHeader)
	if err != nil {
		fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", err.Error(), "\" }"))
		log.Println(w, err.Error())
		return
	}
	for i := range table {
		table[i].convertSize(2)
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
		fmt.Fprintln(w, fmt.Sprint("{\"Error\": \"", err.Error(), "\" }"))
		log.Println(w, err.Error())
		return
	}
	fmt.Fprintln(w, string(bytes))
}

type ViewData struct {
	Separator string
	Folders   []string
}

func createQueryHandler(staticsDir string, startRoot string, htmlPath string) *http.ServeMux {
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
	mux.HandleFunc("/fs", getEntriesTableHandler)

	return mux
}
