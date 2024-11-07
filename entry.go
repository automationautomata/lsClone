package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	// KILOBYTE - количество байт в Килобайте.
	KILOBYTE = 1024

	// MEGABYTE - количество байт в Мегабайте.
	MEGABYTE = 1024 * KILOBYTE

	// GIGABYTE - количество байт в Гигабайте.
	GIGABYTE = 1024 * MEGABYTE
)

// EntryType - тип сущности: папка либо файл
type EntryType string

const (
	File   EntryType = "File"
	Folder EntryType = "Folder"
)

// lsCloneInfo - содержит информацию для вывода на экран.
type lsCloneInfo struct {
	sync.Mutex
	Name          string    `josn:"Name"`
	Type          EntryType `json:"Type"`
	size          int64
	ConvertedSize string `json:"ConvertedSize"`
}

func NewlsCloneInfo(name string, isdir bool) *lsCloneInfo {
	if isdir {
		return &lsCloneInfo{Name: name, Type: Folder}
	}
	return &lsCloneInfo{Name: name, Type: File}
}

// IsDir - проверяет является ли сущность директорией
func (i *lsCloneInfo) IsDir() bool {
	return i.Type == Folder
}

// convertSize - возвращает размер, в зависимости от пересечение границы 1 ГБ / 1 МБ / 1 КБ,
// в виде строки с указанием единиц измерения.
func (i *lsCloneInfo) convertSize(prec int) {
	if i.size >= GIGABYTE {
		i.ConvertedSize = strconv.FormatFloat(float64(i.size)/GIGABYTE, 'f', prec, 64) + " GB"
	}
	if i.size >= MEGABYTE {
		i.ConvertedSize = strconv.FormatFloat(float64(i.size)/MEGABYTE, 'f', prec, 64) + " MB"
	}
	i.ConvertedSize = strconv.FormatFloat(float64(i.size)/KILOBYTE, 'f', prec, 64) + " KB"
}

// IncreaseBy - блокирующее увеличение размера.
func (i *lsCloneInfo) IncreaseBy(size int64) error {
	if i.Type == Folder || i.size == 0 {
		i.Lock()
		i.size += size
		i.Unlock()
		return nil
	}
	return errors.New("The file size is constant")
}

// GetSize - возвращает размер файла в байтах
func (i *lsCloneInfo) GetSize() int64 {
	return i.size
}

// calcSize - позволяет рассчитать размер директории
func (i *lsCloneInfo) calcSize(path string) error {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			err = i.IncreaseBy(info.Size())
			if err != nil {
				return err
			}
			return nil
		})
	if err == nil {
		i.convertSize(2)
	}
	return err
}

// getEntries - возвращает список с инофрмацией о сущностях в директории root
func getEntries(root string) ([]*lsCloneInfo, error) {
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