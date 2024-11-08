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

func checkInput(path string, sort string) (bool, error) {
	if path == "" {
		return false, errors.New("Укажите корневую директорию")

	}
	if _, err := os.Stat(path); err != nil {
		fmt.Println(err.Error())
		return false, errors.New("Корневая директория не существует")
	}

	switch strings.ToLower(sort) {
	case "asc":
		return true, nil
	case "desc":
		return false, nil
	default:
		return false, errors.New("Указан неверный тип сортировки")
	}
}

func getEntriesTableHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "application/json")

	sortHeader := r.Form.Get("sort")
	if sortHeader == "" {
		sortHeader = "asc"
	}
	rootHeader := r.Form.Get("root")

	sortType, err := checkInput(rootHeader, sortHeader)
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

func createQueryHandler(staticsDir string, htmlPath string) *http.ServeMux {
	mux := http.NewServeMux()
	tpl := template.Must(template.ParseFiles(htmlPath))

	fs := http.FileServer(http.Dir(staticsDir))
	staticsPrefix := fmt.Sprint("/", filepath.Base(staticsDir), "/")
	mux.Handle(staticsPrefix, http.StripPrefix(staticsPrefix, fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	})
	mux.HandleFunc("/fs", getEntriesTableHandler)

	return mux
}
