package repositories

import "github.com/kidboy-man/8-level-desent/app/models"

type BookFilter struct {
	Author string
	Page   int
	Limit  int
}

type BookRepository interface {
	Create(book *models.Book) (*models.Book, error)
	FindAll(filter BookFilter) ([]*models.Book, error)
	FindByID(id string) (*models.Book, error)
	Update(id string, book *models.Book) (*models.Book, error)
	Delete(id string) error
	Count(filter BookFilter) (int, error)
}
