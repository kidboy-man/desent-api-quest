package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	httputil "github.com/kidboy-man/8-level-desent/app/controllers/http"
	"github.com/kidboy-man/8-level-desent/app/models"
	"github.com/kidboy-man/8-level-desent/app/repositories"
	"github.com/kidboy-man/8-level-desent/app/services"
)

type BookController struct {
	bookService *services.BookService
}

func NewBookController(bookService *services.BookService) *BookController {
	return &BookController{bookService: bookService}
}

func (ctrl *BookController) Create(c *gin.Context) {
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnError(c, err)
		return
	}

	book, err := ctrl.bookService.CreateBook(&req)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccess(c, http.StatusCreated, book)
}

func (ctrl *BookController) GetAll(c *gin.Context) {
	filter := repositories.BookFilter{
		Author: c.Query("author"),
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	books, total, err := ctrl.bookService.GetAllBooks(filter)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccessWithPagination(c, http.StatusOK, books, total, filter.Page, filter.Limit)
}

func (ctrl *BookController) GetByID(c *gin.Context) {
	id := c.Param("id")

	book, err := ctrl.bookService.GetBookByID(id)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccess(c, http.StatusOK, book)
}

func (ctrl *BookController) Update(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.ReturnError(c, err)
		return
	}

	book, err := ctrl.bookService.UpdateBook(id, &req)
	if err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccess(c, http.StatusOK, book)
}

func (ctrl *BookController) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := ctrl.bookService.DeleteBook(id); err != nil {
		httputil.ReturnError(c, err)
		return
	}

	httputil.ReturnSuccess(c, http.StatusNoContent, nil)
}
