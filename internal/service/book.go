package service

import (
	"context"
	"database/sql"
	"errors"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/yusriltakeuchi/golang-perpustakaan/domain"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/config"
)

type bookService struct {
	cnf                 *config.Config
	bookRepositoy       domain.BookRepository
	bookStockRepository domain.BookStockRepository
	mediaRepository     domain.MediaRepository
}

func NewBook(bookRepository domain.BookRepository,
	bookStockRepository domain.BookStockRepository,
	mediaRepository domain.MediaRepository,
	cnf *config.Config) domain.BookService {
	return &bookService{
		cnf:                 cnf,
		bookRepositoy:       bookRepository,
		bookStockRepository: bookStockRepository,
		mediaRepository:     mediaRepository,
	}
}

// Index implements [domain.BookService].
func (b *bookService) Index(ctx context.Context) ([]dto.BookData, error) {
	result, err := b.bookRepositoy.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	coverId := make([]string, 0)
	for _, v := range result {
		if v.CoverId.Valid {
			coverId = append(coverId, v.CoverId.String)
		}
	}
	covers := make(map[string]string)
	if len(coverId) > 0 {
		coversDb, _ := b.mediaRepository.FindByIds(ctx, coverId)
		for _, v := range coversDb {
			covers[v.Id] = path.Join(b.cnf.Server.Asset, v.Path)
		}
	}

	var data []dto.BookData
	for _, v := range result {
		var coverUrl string
		if v2, e := covers[v.CoverId.String]; e {
			coverUrl = v2
		}

		data = append(data, dto.BookData{
			Id:          v.Id,
			Isbn:        v.Isbn,
			Title:       v.Title,
			CoverUrl:    coverUrl,
			Description: v.Description,
		})
	}
	return data, nil
}

// Show implements [domain.BookService].
func (b *bookService) Show(ctx context.Context, id string) (dto.BookShowData, error) {
	data, err := b.bookRepositoy.FindById(ctx, id)
	if err != nil {
		return dto.BookShowData{}, err
	}
	if data.Id == "" {
		return dto.BookShowData{}, domain.BookNotFound
	}
	stocks, err := b.bookStockRepository.FindByBookId(ctx, data.Id)
	if err != nil {
		return dto.BookShowData{}, err
	}

	stocksData := make([]dto.BookStockData, 0)
	for _, v := range stocks {
		stocksData = append(stocksData, dto.BookStockData{
			Code:   v.Code,
			Status: v.Status,
		})
	}

	var coverUrl string
	if data.CoverId.Valid {
		cover, _ := b.mediaRepository.FindById(ctx, data.CoverId.String)
		if cover.Path != "" {
			coverUrl = path.Join(b.cnf.Server.Asset, cover.Path)
		}
	}

	return dto.BookShowData{
		BookData: dto.BookData{
			Id:          data.Id,
			Isbn:        data.Isbn,
			Title:       data.Title,
			CoverUrl:    coverUrl,
			Description: data.Description,
		},
		Stocks: stocksData,
	}, nil
}

// Create implements [domain.BookService].
func (b *bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}
	book := domain.Book{
		Id:          uuid.NewString(),
		Isbn:        req.Isbn,
		Title:       req.Title,
		Description: req.Description,
		CoverId:     coverId,
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
	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}
	persisted.Isbn = req.Isbn
	persisted.Title = req.Title
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	persisted.CoverId = coverId
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
