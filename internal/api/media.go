package api

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yusriltakeuchi/golang-perpustakaan/domain"
	"github.com/yusriltakeuchi/golang-perpustakaan/dto"
	"github.com/yusriltakeuchi/golang-perpustakaan/internal/config"
)

type mediaApi struct {
	cnf          *config.Config
	mediaService domain.MediaService
}

func NewMedia(app *fiber.App,
	cnf *config.Config,
	mediaService domain.MediaService,
	authMid fiber.Handler) {

	ma := mediaApi{
		cnf:          cnf,
		mediaService: mediaService,
	}
	app.Post("/media", authMid, ma.Create)
	app.Static("/media", cnf.Storage.BasePath)
}

func (ma mediaApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	file, err := ctx.FormFile("media")
	if err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	fileName := uuid.NewString() + "_" + file.Filename
	path := filepath.Join(ma.cnf.Storage.BasePath, fileName)
	err = ctx.SaveFile(file, path)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	res, err := ma.mediaService.Create(c, dto.CreateMediaRequest{
		Path: fileName,
	})
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}
	return ctx.Status(http.StatusCreated).
		JSON(dto.CreateResponseSuccess(res))
}
