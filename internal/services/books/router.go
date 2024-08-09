package books_router

import (
	"noto/internal/services/books/handler"
	"noto/internal/services/books/repository"
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

func BooksRouter(router fiber.Router) {

	repo := repository.NewBookRepository()
	serv := service.NewBookService(repo)
	hand := handler.NewBookHandler(serv)

	router.Post("/books", hand.CreateBook)
	router.Put("/books/:id", hand.UpdateBook)
	router.Get("/books", hand.GetBooks)
	router.Get("/books/:id", hand.GetBook)
	router.Patch("/books/:id", hand.ArchiveBook)
}
