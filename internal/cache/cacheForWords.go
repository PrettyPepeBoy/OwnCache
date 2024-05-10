package cache

import (
	"slices"
	"sync"
	"time"
)

type CacheForWords struct {
	mutex           sync.RWMutex
	cleanupInterval time.Duration
	words           Words
}

type Words struct {
	finish    bool
	wordsUsed int
	letter    byte
	letters   []*Words
}

func NewChForWords(cleanupInterval time.Duration) *CacheForWords {
	return &CacheForWords{
		cleanupInterval: cleanupInterval,
		words:           Words{}}
}

func (c *CacheForWords) Put(letter byte) *Words {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.words.PutLetter(letter)
}

func (w *Words) PutLetter(letter byte) *Words {
	if w.letters == nil {
		w.letters = append(w.letters, &Words{letter: letter})
		return w
	}

	value := slices.Index(w.letters, &Words{letter: letter})

	if value != -1 {
		return w.letters[value].PutLetter(letter)
	}

	w.letters = append(w.letters, &Words{letter: letter})
	return w
}
