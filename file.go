package main

import (
	"strconv"
	"sync"
)

const MEGABYTE = 1048576
const KILOBYTE = 1024
const GIGABYTE = 1073741824

type lsCloneInfo struct {
	Name  string
	IsDir bool
	Size  int64
	sync.Mutex
}

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

func (i *lsCloneInfo) Extend(Size int64) {
	i.Lock()
	i.Size += Size
	i.Unlock()
}
