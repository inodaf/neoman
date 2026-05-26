package models

import "time"

type DocPage struct {
	Author         string
	Repository     string
	RelativePath   string
	Title          string
	Content        string
	LastModifiedAt time.Time
	Vector         map[int]float32
}
