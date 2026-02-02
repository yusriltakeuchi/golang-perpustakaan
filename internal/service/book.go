package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yusriltakeuchi/golang-perpustakaan/domain"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
)

type bookService struct {
	bookRepositoy       domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBook(bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository) domain.BookService {
	return &bookService{
		bookRepositoy:       bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

// Index implements [domain.BookService].
func (b *bookService) Index(ctx context.Context) ([]dto.BookData, error) {
	result, err := b.bookRepositoy.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var data []dto.BookData
	for _, v := range result {
		data = append(data, dto.BookData{
			Id:          v.Id,
			Isbn:        v.Isbn,
			Title:       v.Title,
			Description: v.Description,
		})
	}
	return data, nil
}

// Show implements [domain.BookService].
func (b *bookService) Show(ctx context.Context, id string) (dto.BookData, error) {
	data, err := b.bookRepositoy.FindById(ctx, id)
	if err != nil {
		return dto.BookData{}, err
	}
	if data.Id == "" {
		return dto.BookData{}, errors.New("data buku tidak ditemukan")
	}
	return dto.BookData{
		Id:          data.Id,
		Isbn:        data.Isbn,
		Title:       data.Title,
		Description: data.Description,
	}, nil
}

// Create implements [domain.BookService].
func (b *bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	book := domain.Book{
		Id:          uuid.NewString(),
		Isbn:        req.Isbn,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   sql.NullTime{Valid: true, Time: time.Now()},
	}
	return b.bookRepositoy.Save(ctx, &book)
}

// Update implements [domain.BookService].
func (b *bookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	persisted, err := b.bookRepositoy.FindById(ctx, req.Id)
	if err != nil {
		return err
	}
	if persisted.Id == "" {
		return errors.New("data buku tidak ditemukan")
	}
	persisted.Isbn = req.Isbn
	persisted.Title = req.Title
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	return b.bookRepositoy.Update(ctx, &persisted)
}

// Delete implements [domain.BookService].
func (b *bookService) Delete(ctx context.Context, id string) error {
	persisted, err := b.bookRepositoy.FindById(ctx, id)
	if err != nil {
		return err
	}
	if persisted.Id == "" {
		return errors.New("data buku tidak ditemukan")
	}
	err = b.bookRepositoy.Delete(ctx, persisted.Id)
	if err != nil {
		return err
	}
	return b.bookStockRepository.DeleteByBookId(ctx, persisted.Id)
}
