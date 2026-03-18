package storage

import (
	"errors"
	"fmt"
	"url-shortener/models"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

type URLEntry = models.URLEntry

type MemoryStorage struct {
	entries map[string]*URLEntry
}

type Storage interface {
	Save(entry *URLEntry) error
	Get(code string) (*URLEntry, error)
	All() []*URLEntry
}

func (m *MemoryStorage) Save(entry *URLEntry) error {
	if _, exists := m.entries[entry.Code]; exists {
		return fmt.Errorf("storage.save %q: %w", entry.Code, ErrAlreadyExists)
	}

	m.entries[entry.Code] = entry

	return nil
}

func (m *MemoryStorage) Get(code string) (*URLEntry, error) {
	val, ok := m.entries[code]

	if !ok {
		return nil, fmt.Errorf("storage.get %q: %w", code, ErrNotFound)
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
