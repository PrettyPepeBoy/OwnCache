package cache

import (
	"log/slog"
	"sync"
	"time"
)

type WordsCache struct {
	mutex           sync.RWMutex
	cleanupInterval time.Duration
	words           Word
}

type Word struct {
	letters    []*Word
	wordsCount int
	finish     bool
	letter     byte
}

func NewWordsCache(cleanupInterval time.Duration, logger *slog.Logger) *WordsCache {
	words := InitWord('a')
	cache := WordsCache{words: words, cleanupInterval: cleanupInterval}
	cache.StartGC(logger)
	return &cache
}

func InitWord(letter byte) Word {
	var words []*Word
	w := Word{
		letters: append(words, &Word{letter: letter}),
	}
	return w
}

func (w *Word) PutLetters(letter []byte) *Word {
	w.wordsCount++
	if len(letter) == 0 {
		w.finish = true
		return w
	}

	//TODO добавить валидацию слова
	if w.letters == nil {
		w.letters = append(w.letters, &Word{letter: letter[0]})
		return w.letters[0].PutLetters(letter[1:])
	}

	found, index := binarySearch(w.letters, letter[0])

	if !found {
		w.letters = append(w.letters, &Word{letter: letter[0]})
		binarySort(w.letters, 0, len(w.letters)-1, partition)
		_, index = binarySearch(w.letters, letter[0])
		return w.letters[index].PutLetters(letter[1:])
	}
	return w.letters[index].PutLetters(letter[1:])
}

func (w *Word) GetWord(word []byte) []byte {
	node, found := w.countWord(word)
	if !found {
		return nil
	}
	var newWord []byte
	newWord = *node.maxCount(&newWord)
	word = append(word, newWord...)
	return word
}

func (w *Word) countWord(word []byte) (*Word, bool) {
	if len(word) == 0 {
		return w, true
	}
	found, index := binarySearch(w.letters, word[0])
	if !found {
		return nil, false
	}
	return w.letters[index].countWord(word[1:])
}

func (w *Word) maxCount(word *[]byte) *[]byte {
	if w.finish {
		return word
	}
	var count, index int

	for i, elem := range w.letters {
		if count < elem.wordsCount {
			count = elem.wordsCount
			index = i
		}
	}

	*word = append(*word, w.letters[index].letter)

	return w.letters[index].maxCount(word)

}

func binarySearch(letters []*Word, letter byte) (bool, int) {
	length := len(letters)
	if letters[length/2].letter == letter {
		return true, length / 2
	}
	if length == 1 {
		return false, -1
	}

	if letters[length/2].letter > letter {
		return binarySearch(letters[:length/2], letter)
	}

	return binarySearch(letters[length/2:], letter)
}

// todo бинарная вставка
func binarySort(arr []*Word, low int, high int, fn func([]*Word, int, int) int) {
	if low < high {
		index := fn(arr, low, high)
		binarySort(arr, low, index-1, fn)
		binarySort(arr, index+1, high, fn)
	}

}

func partition(arr []*Word, low int, high int) int {
	i := low - 1
	pivot := arr[high].letter
	for j := low; j < high; j++ {
		if arr[j].letter < pivot {
			i++

			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[high], arr[i+1] = arr[i+1], arr[high]
	return i + 1
}

func (wc *WordsCache) PutWord(letters []byte) *Word {
	return wc.words.PutLetters(letters)
}

func (wc *WordsCache) GetWord(letters []byte) []byte {
	return wc.words.GetWord(letters)
}

func (wc *WordsCache) StartGC(logger *slog.Logger) {
	go wc.GC(logger)
}

func (wc *WordsCache) GC(logger *slog.Logger) {
	for {
		time.Sleep(wc.cleanupInterval)
		logger.Info("starting clear cache")

		wc.words = InitWord('a')
	}
}
