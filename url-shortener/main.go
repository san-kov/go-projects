package main

import (
	"errors"
	"fmt"
	"url-shortener/models"
	"url-shortener/storage"
)

func main() {
	e1 := models.NewURLEntry("abc1", "https://google.com")
	e2 := models.NewURLEntry("abc2", "https://twitter.com")
	e3 := models.NewURLEntry("abc3", "https://youtube.com")

	s := storage.NewMemoryStorage()

	s.Save(e1)
	s.Save(e2)
	s.Save(e3)

	var _ storage.Storage = (*storage.MemoryStorage)(nil)

	err := s.Save(e3)

	if errors.Is(err, storage.ErrAlreadyExists) {
		fmt.Println("already exists")
	}

	_, err2 := s.Get("qwe")
	if errors.Is(err2, storage.ErrNotFound) {
		fmt.Println("not found")
	}

	printAll(s)
}

func printAll(s storage.Storage) {
	entries := s.All()

	for _, entry := range entries {
		fmt.Println(entry.String())
	}

}
