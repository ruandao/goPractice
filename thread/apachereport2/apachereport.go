package main

import "sync"

type pageMap struct {
	countForPage map[string]int
	mutex *sync.RWMutex
}

func NewPageMap() *pageMap {
	return &pageMap{make(map[string]int), new(sync.RWMutex)}
}

func (pm *pageMap)Increment(page string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.countForPage[page]++
}