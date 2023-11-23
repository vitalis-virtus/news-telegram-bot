package model

import "time"

type Item struct { //for RSS
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourceName string
}

type Source struct {
	ID        int
	Name      string
	FeedURL   string
	CreatedAt time.Time
}

type Article struct { //article in the format in our system
	ID          int
	Source      string
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	CreatedAt   time.Time
	PostedAt    time.Time
}
