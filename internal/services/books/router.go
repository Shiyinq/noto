package books_router

import (
	"noto/internal/config"
	"noto/internal/services/books/handler"
	"noto/internal/services/books/repository"
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

func BooksRouter(router fiber.Router) {

	repo := repository.NewBookRepository(config.DB)
	serv := service.NewBookService(repo)
	hand := handler.NewBookHandler(serv)

	router.Post("/books", hand.CreateBook)
	router.Put("/books/:bookId", hand.UpdateBook)
	router.Get("/books", hand.GetBooks)
	router.Get("/books/:bookId", hand.GetBook)
	router.Patch("/books/:bookId", hand.ArchiveBook)
}
