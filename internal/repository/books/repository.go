package books

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ferdy-adr/elibrary-backend/internal/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateBook(book *model.Book) error {
	query := `
		INSERT INTO books (title, isbn, year, publisher, author, cover_image, synopsis) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, book.Title, book.ISBN, book.Year, book.Publisher, book.Author, book.CoverImage, book.Synopsis)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	book.ID = int(id)
	return nil
}

func (r *Repository) GetBookByID(id int) (*model.Book, error) {
	book := &model.Book{}
	query := `
		SELECT id, title, isbn, year, publisher, author, cover_image, synopsis, created_at, updated_at 
		FROM books 
		WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.Title, &book.ISBN, &book.Year, &book.Publisher,
		&book.Author, &book.CoverImage, &book.Synopsis, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r *Repository) GetBooks(params model.BookQueryParams) ([]model.Book, int, error) {
	var books []model.Book
	var total int

	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}

	if params.Search != "" {
		whereConditions = append(whereConditions, "(title LIKE ? OR author LIKE ? OR publisher LIKE ?)")
		searchTerm := "%" + params.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if params.Year > 0 {
		whereConditions = append(whereConditions, "year = ?")
		args = append(args, params.Year)
	}

	if params.Publisher != "" {
		whereConditions = append(whereConditions, "publisher LIKE ?")
		args = append(args, "%"+params.Publisher+"%")
	}

	if params.Author != "" {
		whereConditions = append(whereConditions, "author LIKE ?")
		args = append(args, "%"+params.Author+"%")
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM books %s", whereClause)
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (params.Page - 1) * params.Limit
	query := fmt.Sprintf(`
		SELECT id, title, isbn, year, publisher, author, cover_image, synopsis, created_at, updated_at 
		FROM books %s 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, params.Limit, offset)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var book model.Book
		err := rows.Scan(
			&book.ID, &book.Title, &book.ISBN, &book.Year, &book.Publisher,
			&book.Author, &book.CoverImage, &book.Synopsis, &book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, book)
	}

	return books, total, nil
}

func (r *Repository) UpdateBook(id int, book *model.Book) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}

	if book.Title != "" {
		setParts = append(setParts, "title = ?")
		args = append(args, book.Title)
	}

	if book.ISBN != "" {
		setParts = append(setParts, "isbn = ?")
		args = append(args, book.ISBN)
	}

	if book.Year > 0 {
		setParts = append(setParts, "year = ?")
		args = append(args, book.Year)
	}

	if book.Publisher != "" {
		setParts = append(setParts, "publisher = ?")
		args = append(args, book.Publisher)
	}

	if book.Author != "" {
		setParts = append(setParts, "author = ?")
		args = append(args, book.Author)
	}

	if book.CoverImage != "" {
		setParts = append(setParts, "cover_image = ?")
		args = append(args, book.CoverImage)
	}

	if book.Synopsis != "" {
		setParts = append(setParts, "synopsis = ?")
		args = append(args, book.Synopsis)
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf("UPDATE books SET %s WHERE id = ?", strings.Join(setParts, ", "))
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *Repository) DeleteBook(id int) error {
	query := "DELETE FROM books WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *Repository) CheckISBNExists(isbn string, excludeID int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM books WHERE isbn = ? AND id != ?"
	err := r.db.QueryRow(query, isbn, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
