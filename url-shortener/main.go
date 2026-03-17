package main

import (
	"errors"
	"fmt"
	"time"
)

type MemoryStorage struct {
	entries map[string]*URLEntry
}

type Storage interface {
	Save(entry *URLEntry) error
	Get(code string) (*URLEntry, error)
	All() []*URLEntry
}

type URLEntry struct {
	Code        string
	OriginalURL string
	Clicks      int
	CreatedAt   time.Time
}

func (m *MemoryStorage) Save(entry *URLEntry) error {
	m.entries[entry.Code] = entry

	return nil
}

func (m *MemoryStorage) Get(code string) (*URLEntry, error) {
	val, ok := m.entries[code]

	if !ok {
		return nil, errors.New("Cannot get value")
	}

	return val, nil
}

func (m *MemoryStorage) All() []*URLEntry {
	var entries []*URLEntry

	for _, entry := range m.entries {
		entries = append(entries, entry)
	}

	return entries
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries: make(map[string]*URLEntry),
	}
}

func main() {
	e1 := NewURLEntry("abc1", "https://google.com")
	e2 := NewURLEntry("abc2", "https://twitter.com")
	e3 := NewURLEntry("abc3", "https://youtube.com")

	storage := NewMemoryStorage()

	storage.Save(e1)
	storage.Save(e2)
	storage.Save(e3)

	var _ Storage = (*MemoryStorage)(nil)

	printAll(storage)
}

func NewURLEntry(code, url string) *URLEntry {
	return &URLEntry{
		Code:        code,
		OriginalURL: url,
		Clicks:      0,
		CreatedAt:   time.Now(),
	}

}

func (e *URLEntry) Click() {
	e.Clicks++
}

func (e *URLEntry) String() string {
	return fmt.Sprintf("[%s] %s (%d Clicks)", e.Code, e.OriginalURL, e.Clicks)
}

func printAll(s Storage) {
	entries := s.All()

	for _, entry := range entries {
		fmt.Println(entry.String())
	}

}
