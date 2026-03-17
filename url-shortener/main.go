package main

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

func main() {
	e1 := NewURLEntry("abc1", "https://google.com")
	e2 := NewURLEntry("abc2", "https://twitter.com")

	e1.Click()
	e1.Click()
	e1.Click()

	e2.Click()
	e2.Click()
	e2.Click()

	fmt.Println(e1.String())
	fmt.Println(e2.String())
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
