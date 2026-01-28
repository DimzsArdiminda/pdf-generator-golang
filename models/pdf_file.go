package models

import "time"

type PdfFile struct {
	ID           int64      `json:"id"`
	Filename     string     `json:"filename"`
	OriginalName *string    `json:"original_name"`
	Filepath     string     `json:"filepath"`
	Size         *int64     `json:"size"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
