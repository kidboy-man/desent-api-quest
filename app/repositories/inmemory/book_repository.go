package inmemory

import (
	"strings"
	"sync"

	apperrors "github.com/kidboy-man/8-level-desent/app/errors"
	"github.com/kidboy-man/8-level-desent/app/models"
	"github.com/kidboy-man/8-level-desent/app/repositories"
)

type BookRepository struct {
	mu    sync.RWMutex
	store map[string]*models.Book
	order []string // maintains insertion order
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		store: make(map[string]*models.Book),
		order: make([]string, 0),
	}
}

func (r *BookRepository) Create(book *models.Book) (*models.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored := *book
	r.store[book.ID] = &stored
	r.order = append(r.order, book.ID)

	result := stored
	return &result, nil
}

func (r *BookRepository) FindAll(filter repositories.BookFilter) ([]*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var filtered []*models.Book
	for _, id := range r.order {
		book, exists := r.store[id]
		if !exists {
			continue
		}
		if filter.Author != "" && !strings.EqualFold(book.Author, filter.Author) {
			continue
		}
		copy := *book
		filtered = append(filtered, &copy)
	}

	if filter.Page > 0 && filter.Limit > 0 {
		start := (filter.Page - 1) * filter.Limit
		if start >= len(filtered) {
			return []*models.Book{}, nil
		}
		end := start + filter.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
		filtered = filtered[start:end]
	}

	return filtered, nil
}

func (r *BookRepository) FindByID(id string) (*models.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, exists := r.store[id]
	if !exists {
		return nil, apperrors.Newf(404, "NOT_FOUND", "book with id %s not found", id)
	}

	result := *book
	return &result, nil
}

func (r *BookRepository) Update(id string, book *models.Book) (*models.Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return nil, apperrors.Newf(404, "NOT_FOUND", "book with id %s not found", id)
	}

	book.ID = id
	stored := *book
	r.store[id] = &stored

	result := stored
	return &result, nil
}

func (r *BookRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return apperrors.Newf(404, "NOT_FOUND", "book with id %s not found", id)
	}

	delete(r.store, id)
	for i, oid := range r.order {
		if oid == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}

	return nil
}

func (r *BookRepository) Count(filter repositories.BookFilter) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, id := range r.order {
		book, exists := r.store[id]
		if !exists {
			continue
		}
		if filter.Author != "" && !strings.EqualFold(book.Author, filter.Author) {
			continue
		}
		count++
	}

	return count, nil
}
