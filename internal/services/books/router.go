package books_router

import (
	"noto/internal/services/books/handler"
	"noto/internal/services/books/repository"
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

var repo = repository.NewBookRepository()
var serv = service.NewBookService(repo)
var hand = handler.NewBookHandler(serv)

func BooksRouter(app *fiber.App) {
	app.Get("/books", hand.GetBooks)
	app.Get("/books/:id", hand.GetBook)
}
