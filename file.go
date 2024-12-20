package main

import (
	"errors"
	"strconv"
	"sync"
)

// KILOBYTE - количество байт в Килобайте.
const KILOBYTE = 1024

// MEGABYTE - количество байт в Мегабайте.
const MEGABYTE = 1024 * KILOBYTE

// GIGABYTE - количество байт в Гигабайте.
const GIGABYTE = 1024 * MEGABYTE

// lsCloneInfo - содержит информацию для вывода на экран.
type lsCloneInfo struct {
	Name  string
	IsDir bool
	size  int64
	sync.Mutex
}

// convertSize - возвращает размер, в зависимости от пересечение границы 1 ГБ / 1 МБ / 1 КБ,
// в виде строки с указанием единиц измерения.
func (i *lsCloneInfo) convertSize(prec int) string {
	if i.size >= GIGABYTE {
		return strconv.FormatFloat(float64(i.size)/GIGABYTE, 'f', prec, 64) + " GB"
	}
	if i.size >= MEGABYTE {
		return strconv.FormatFloat(float64(i.size)/MEGABYTE, 'f', prec, 64) + " MB"
	}
	return strconv.FormatFloat(float64(i.size)/KILOBYTE, 'f', prec, 64) + " KB"
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

// Получить размер файла
func (i *lsCloneInfo) GetSize() int64 {
	return i.size
}
