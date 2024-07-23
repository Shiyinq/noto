package books_router

import (
	"noto/internal/services/books/handler"
	"noto/internal/services/books/repository"
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

func BooksRouter(app *fiber.App) {

	repo := repository.NewBookRepository()
	serv := service.NewBookService(repo)
	hand := handler.NewBookHandler(serv)

	app.Post("/books", hand.CreateBook)
	app.Put("/books/:id", hand.UpdateBook)
	app.Get("/books", hand.GetBooks)
	app.Get("/books/:id", hand.GetBook)
}
