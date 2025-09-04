package books

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ferdy-adr/elibrary-backend/internal/configs"
	"github.com/ferdy-adr/elibrary-backend/internal/model"
	bookRepo "github.com/ferdy-adr/elibrary-backend/internal/repository/books"
)

type Service struct {
	bookRepository *bookRepo.Repository
}

func NewService(bookRepository *bookRepo.Repository) *Service {
	return &Service{
		bookRepository: bookRepository,
	}
}

func (s *Service) CreateBook(req model.CreateBookRequest, coverFile *multipart.FileHeader) (*model.Book, error) {
	// Check if ISBN already exists
	exists, err := s.bookRepository.CheckISBNExists(req.ISBN, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("ISBN already exists")
	}

	book := &model.Book{
		Title:     req.Title,
		ISBN:      req.ISBN,
		Year:      req.Year,
		Publisher: req.Publisher,
		Author:    req.Author,
		Synopsis:  req.Synopsis,
	}

	// Handle cover image upload if provided
	if coverFile != nil {
		coverImagePath, err := s.uploadCoverImage(coverFile)
		if err != nil {
			return nil, fmt.Errorf("failed to upload cover image: %v", err)
		}
		book.CoverImage = coverImagePath
	}

	err = s.bookRepository.CreateBook(book)
	if err != nil {
		// If book creation fails and we uploaded an image, clean it up
		if book.CoverImage != "" {
			s.deleteCoverImage(book.CoverImage)
		}
		return nil, err
	}

	return book, nil
}

func (s *Service) GetBooks(params model.BookQueryParams) (*model.BookListResponse, error) {
	// Set default values
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	books, total, err := s.bookRepository.GetBooks(params)
	if err != nil {
		return nil, err
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &model.BookListResponse{
		Books:      books,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *Service) GetBookByID(id int) (*model.Book, error) {
	return s.bookRepository.GetBookByID(id)
}

func (s *Service) UpdateBook(id int, req model.UpdateBookRequest, coverFile *multipart.FileHeader) (*model.Book, error) {
	// Check if book exists
	existingBook, err := s.bookRepository.GetBookByID(id)
	if err != nil {
		return nil, errors.New("book not found")
	}

	// Check ISBN uniqueness if ISBN is being updated
	if req.ISBN != "" && req.ISBN != existingBook.ISBN {
		exists, err := s.bookRepository.CheckISBNExists(req.ISBN, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("ISBN already exists")
		}
	}

	book := &model.Book{
		Title:     req.Title,
		ISBN:      req.ISBN,
		Year:      req.Year,
		Publisher: req.Publisher,
		Author:    req.Author,
		Synopsis:  req.Synopsis,
	}

	// Handle cover image upload if provided
	if coverFile != nil {
		coverImagePath, err := s.uploadCoverImage(coverFile)
		if err != nil {
			return nil, fmt.Errorf("failed to upload cover image: %v", err)
		}

		// Delete old cover image if it exists
		if existingBook.CoverImage != "" {
			s.deleteCoverImage(existingBook.CoverImage)
		}

		book.CoverImage = coverImagePath
	}

	err = s.bookRepository.UpdateBook(id, book)
	if err != nil {
		// If update fails and we uploaded a new image, clean it up
		if book.CoverImage != "" && book.CoverImage != existingBook.CoverImage {
			s.deleteCoverImage(book.CoverImage)
		}
		return nil, err
	}

	// Get updated book
	return s.bookRepository.GetBookByID(id)
}

func (s *Service) DeleteBook(id int) error {
	// Get book to check if it has a cover image
	book, err := s.bookRepository.GetBookByID(id)
	if err != nil {
		return errors.New("book not found")
	}

	// Delete from database
	err = s.bookRepository.DeleteBook(id)
	if err != nil {
		return err
	}

	// Delete cover image if exists
	if book.CoverImage != "" {
		s.deleteCoverImage(book.CoverImage)
	}

	return nil
}

func (s *Service) uploadCoverImage(fileHeader *multipart.FileHeader) (string, error) {
	// Validate file type
	if !s.isValidImageType(fileHeader.Filename) {
		return "", errors.New("invalid file type. Only JPG, JPEG, PNG files are allowed")
	}

	// Create upload directory if it doesn't exist
	uploadPath := configs.Get().Upload.Path
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("cover_%d%s", time.Now().UnixNano(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	// Return relative path for storing in database
	return fmt.Sprintf("/images/%s", filename), nil
}

func (s *Service) deleteCoverImage(imagePath string) {
	if imagePath == "" {
		return
	}

	// Extract filename from path
	filename := filepath.Base(imagePath)
	fullPath := filepath.Join(configs.Get().Upload.Path, filename)

	// Delete file (ignore errors)
	os.Remove(fullPath)
}

func (s *Service) isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExtensions := []string{".jpg", ".jpeg", ".png"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}
