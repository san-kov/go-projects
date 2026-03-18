package storage

import (
	"errors"
	"fmt"
	"url-shortener/models"
)

var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

var _ Storage = (*MemoryStorage)(nil)

type MemoryStorage struct {
	entries map[string]*models.URLEntry
}

type Storage interface {
	Save(entry *models.URLEntry) error
	Get(code string) (*models.URLEntry, error)
	All() []*models.URLEntry
}

func (m *MemoryStorage) Save(entry *models.URLEntry) error {
	if _, exists := m.entries[entry.Code]; exists {
		return fmt.Errorf("storage.save %q: %w", entry.Code, ErrAlreadyExists)
	}

	m.entries[entry.Code] = entry

	return nil
}

func (m *MemoryStorage) Get(code string) (*models.URLEntry, error) {
	val, ok := m.entries[code]

	if !ok {
		return nil, fmt.Errorf("storage.get %q: %w", code, ErrNotFound)
	}

	return val, nil
}

func (m *MemoryStorage) All() []*models.URLEntry {
	var entries []*models.URLEntry

	for _, entry := range m.entries {
		entries = append(entries, entry)
	}

	return entries
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries: make(map[string]*models.URLEntry),
	}
}
