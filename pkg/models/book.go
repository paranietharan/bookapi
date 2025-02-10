package models

import "time"

type Book struct {
	ID             string
	Name           string
	Author         string
	Category       string
	Description    string
	LastAccessTime time.Time
}
