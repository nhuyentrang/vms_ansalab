package controllers

import (
	"container/list"
	"sync"
)

// SafeOrderedMap wraps a map, a mutex, and a linked list to provide thread-safe access and order tracking
type SafeOrderedMap struct {
	mu    sync.Mutex
	m     map[interface{}]interface{}
	order *list.List
}

// NewSafeOrderedMap creates a new SafeOrderedMap
func NewSafeOrderedMap() *SafeOrderedMap {
	return &SafeOrderedMap{
		m:     make(map[interface{}]interface{}),
		order: list.New(),
	}
}

// Store adds an element to the map and tracks its order
func (som *SafeOrderedMap) Store(key interface{}, value interface{}) {
	som.mu.Lock()
	defer som.mu.Unlock()
	if _, exists := som.m[key]; !exists {
		som.order.PushBack(key)
	}
	som.m[key] = value
}

// Load retrieves and removes an element from the map
func (som *SafeOrderedMap) Load(key interface{}) (interface{}, bool) {
	som.mu.Lock()
	defer som.mu.Unlock()
	value, ok := som.m[key]
	if ok {
		delete(som.m, key)
		for e := som.order.Front(); e != nil; e = e.Next() {
			if e.Value == key {
				som.order.Remove(e)
				break
			}
		}
	}
	return value, ok
}

// Delete removes an element from the map and updates the order
func (som *SafeOrderedMap) Delete(key interface{}) {
	som.mu.Lock()
	defer som.mu.Unlock()
	if _, exists := som.m[key]; exists {
		delete(som.m, key)
		for e := som.order.Front(); e != nil; e = e.Next() {
			if e.Value == key {
				som.order.Remove(e)
				break
			}
		}
	}
}

// DeleteOldest removes the oldest element from the map
func (som *SafeOrderedMap) DeleteOldest() {
	som.mu.Lock()
	defer som.mu.Unlock()
	if som.order.Len() > 0 {
		oldest := som.order.Front()
		key := oldest.Value
		delete(som.m, key)
		som.order.Remove(oldest)
	}
}

// Count returns the number of elements in the map
func (som *SafeOrderedMap) Count() int {
	som.mu.Lock()
	defer som.mu.Unlock()
	return len(som.m)
}
