package main

import (
	"errors"
	"strconv"
	"sync"
)

// KILOBYTE - количество байт в Килобайте.
const KILOBYTE = 1024

// MEGABYTE - количество байт в Мегабайте.
const MEGABYTE = 1048576

// GIGABYTE - количество байт в Гигабайте.
const GIGABYTE = 1073741824

// lsCloneInfo - содержит информацию для вывода на экран.
type lsCloneInfo struct {
	Name  string
	IsDir bool
	size  int64
	sync.Mutex
}

// convertsize - возвращает размер, в зависимости от пересечение границы 1 ГБ / 1 МБ / 1 КБ,
// в виде строки с указанием единиц измерения.
func (i *lsCloneInfo) convertSize(prec int) string {
	if i.size >= GIGABYTE {
		return strconv.FormatFloat(float64(i.size/GIGABYTE), 'f', prec, 64) + " GB"
	}
	if i.size >= MEGABYTE {
		return strconv.FormatFloat(float64(i.size/MEGABYTE), 'f', prec, 64) + " MB"
	} else {
		return strconv.FormatFloat(float64(i.size)/KILOBYTE, 'f', prec, 64) + " KB"
	}
	//return strconv.Itoa(int(i.size))
}

// IncreaseBy - блокирующее увеличение размера.
func (i *lsCloneInfo) IncreaseBy(size int64) error {
	if i.IsDir || i.size == 0 {
		i.Lock()
		i.size += size
		i.Unlock()
		return nil
	}
	return errors.New("The file size is constant")
}

func (i *lsCloneInfo) GetSize() int64 {
	return i.size
}
