package model

import "time"

type Book struct {
	ID         int       `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	ISBN       string    `json:"isbn" db:"isbn"`
	Year       int       `json:"year" db:"year"`
	Publisher  string    `json:"publisher" db:"publisher"`
	Author     string    `json:"author" db:"author"`
	CoverImage string    `json:"cover_image" db:"cover_image"`
	Synopsis   string    `json:"synopsis" db:"synopsis"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBookRequest struct {
	Title     string `form:"title" binding:"required"`
	ISBN      string `form:"isbn" binding:"required"`
	Year      int    `form:"year" binding:"required"`
	Publisher string `form:"publisher" binding:"required"`
	Author    string `form:"author" binding:"required"`
	Synopsis  string `form:"synopsis"`
}

type UpdateBookRequest struct {
	Title     string `form:"title"`
	ISBN      string `form:"isbn"`
	Year      int    `form:"year"`
	Publisher string `form:"publisher"`
	Author    string `form:"author"`
	Synopsis  string `form:"synopsis"`
}

type BookListResponse struct {
	Books      []Book `json:"books"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
}

type BookQueryParams struct {
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
	Search    string `form:"search"`
	Year      int    `form:"year"`
	Publisher string `form:"publisher"`
	Author    string `form:"author"`
}
