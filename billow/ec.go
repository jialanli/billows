package billow

import (
	"math/rand"
	"sync"
)

var mu sync.Mutex
var idStorage map[int]bool

func initIdStorage(size int) {
	idStorage = make(map[int]bool, size)
}

func getI() int {
	mu.Lock()
	n := len(idStorage)
	for _, v := range idStorage {
		if !v {
			n++
		}
	}
	mu.Unlock()
	return n
}

func dressI(n int, u bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := idStorage[n]; ok {
		idStorage[n] = u
	}
}

func freedI() {
	mu.Lock()
	defer mu.Unlock()
	for k := range idStorage {
		delete(idStorage, k)
	}
}

func generateI(size int) int {
	n, m := size*5, 0
	mu.Lock()
	defer mu.Unlock()
	for !false {
		m = rand.Intn(n)
		if _, ok := idStorage[m]; ok {
			continue
		}
		idStorage[m] = true
		break
	}
	return m
}
