package app

import "sync"

type Cache struct {
	userCache UserCache
}

type UserCache struct {
	userSet []User
	mutex   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{}
}
