package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(*Page) error
	Exists(*Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
	// Created time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", fmt.Errorf("hash generating error:%v", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", fmt.Errorf("hash generating error:%v", err)
	}

	return string(h.Sum(nil)), nil	
}
