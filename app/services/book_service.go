package services

import (
	"strings"

	apperrors "github.com/kidboy-man/8-level-desent/app/errors"
	"github.com/google/uuid"
	"github.com/kidboy-man/8-level-desent/app/models"
	"github.com/kidboy-man/8-level-desent/app/repositories"
)

type BookService struct {
	repo repositories.BookRepository
}

func NewBookService(repo repositories.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(req *models.CreateBookRequest) (*models.Book, error) {
	if err := validateBookFields(req.Title, req.Author, req.Year); err != nil {
		return nil, err
	}

	book := &models.Book{
		ID:     uuid.New().String(),
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
	}
	return s.repo.Create(book)
}

func (s *BookService) GetAllBooks(filter repositories.BookFilter) ([]*models.Book, int, error) {
	countFilter := repositories.BookFilter{Author: filter.Author}
	total, err := s.repo.Count(countFilter)
	if err != nil {
		return nil, 0, err
	}

	books, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (s *BookService) GetBookByID(id string) (*models.Book, error) {
	return s.repo.FindByID(id)
}

func (s *BookService) UpdateBook(id string, req *models.UpdateBookRequest) (*models.Book, error) {
	if err := validateBookFields(req.Title, req.Author, req.Year); err != nil {
		return nil, err
	}

	book := &models.Book{
		Title:  req.Title,
		Author: req.Author,
		Year:   req.Year,
	}
	return s.repo.Update(id, book)
}

func (s *BookService) DeleteBook(id string) error {
	return s.repo.Delete(id)
}

func validateBookFields(title, author string, year int) error {
	var missing []string
	if title == "" {
		missing = append(missing, "title")
	}
	if author == "" {
		missing = append(missing, "author")
	}
	if year <= 0 {
		missing = append(missing, "year")
	}
	if len(missing) > 0 {
		return apperrors.NewBadRequest("missing or invalid fields: " + strings.Join(missing, ", "))
	}
	return nil
}
