package cache

import (
	"log/slog"
	"slices"
	"sync"
	"time"
)

type Cache struct {
	mutex           sync.RWMutex
	cleanupInterval time.Duration
	items           map[int64]Item
}

type Item struct {
	values []int64
}

func NewCache(cleanupInterval time.Duration, logger *slog.Logger) *Cache {
	items := make(map[int64]Item)
	cache := Cache{
		items:           items,
		cleanupInterval: cleanupInterval,
	}
	cache.StartGC(logger)
	return &cache
}

func (c *Cache) SetKey(key, value int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !slices.Contains(c.items[key].values, value) {
		c.items[key] = Item{
			values: append(c.items[key].values, value),
		}
	}
}

func (c *Cache) ShowCache(key int64) ([]int64, bool) {
	value, ok := c.items[key]
	if !ok {
		return nil, false
	}
	return value.values, ok
}

func (c *Cache) StartGC(logger *slog.Logger) {
	go c.GC(logger)
}

func (c *Cache) GC(logger *slog.Logger) {
	for {
		time.Sleep(c.cleanupInterval)
		logger.Info("starting clear cache")

		if c.items == nil {
			return
		}

		c.items = make(map[int64]Item)
	}
}
