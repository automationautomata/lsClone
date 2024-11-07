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

func getInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sortHeader := r.Form.Get("sort")
	if sortHeader == "" {
		sortHeader = "asc"
	}
	rootHeader := r.Form.Get("root")

	sortType, err := checkInput(rootHeader, sortHeader)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		log.Fatalln(w, err.Error())
		return
	}

	table, err := getEntries(rootHeader)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		log.Fatalln(w, err.Error())
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
		fmt.Fprintln(w, err.Error())
		log.Fatalln(w, err.Error())
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
	return mux
}
