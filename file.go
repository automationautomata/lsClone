package main

import (
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
	Size  int64
	sync.Mutex
}

// convertSize - возвращает размер, в зависимости от пересечение границы 1 ГБ / 1 МБ / 1 КБ,
// в виде строки с указанием единиц измерения.
func (i *lsCloneInfo) convertSize(prec int) string {
	if i.Size >= GIGABYTE {
		return strconv.FormatFloat(float64(i.Size/GIGABYTE), 'f', prec, 64) + " GB"
	}
	if i.Size >= MEGABYTE {
		return strconv.FormatFloat(float64(i.Size/MEGABYTE), 'f', prec, 64) + " MB"
	} else {
		return strconv.FormatFloat(float64(i.Size)/KILOBYTE, 'f', prec, 64) + " KB"
	}
	//return strconv.Itoa(int(i.Size))
}

// IncreaseBy - блокирующее увеличение размера.
func (i *lsCloneInfo) IncreaseBy(Size int64) {
	i.Lock()
	i.Size += Size
	i.Unlock()
}
