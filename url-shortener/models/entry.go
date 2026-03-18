package models

import (
	"fmt"
	"time"
)

type URLEntry struct {
	Code        string
	OriginalURL string
	Clicks      int
	CreatedAt   time.Time
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
