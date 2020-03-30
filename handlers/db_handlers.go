package handlers

import (
	"context"
	"database/sql"
	"github.com/pytorchtw/go-db-services-starter/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type DBHandler struct {
	db  *sql.DB
	ctx context.Context
}

func NewDBHandler(db *sql.DB) (*DBHandler, error) {
	h := DBHandler{}
	boil.SetDB(db)
	h.db = db
	h.ctx = context.Background()
	return &h, nil
}

func (h *DBHandler) GetPage(url string) (*models.Page, error) {
	page, err := models.Pages(models.PageWhere.URL.EQ(url)).One(h.ctx, h.db)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (h *DBHandler) CreatePage(page *models.Page) error {
	err := page.Insert(h.ctx, h.db, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (h *DBHandler) DeletePage(url string) (int64, error) {
	page, err := h.GetPage(url)
	rowsAff, err := page.Delete(h.ctx, h.db)
	if err != nil {
		return 0, err
	}
	return rowsAff, nil
}

func (h *DBHandler) Close() error {
	err := h.db.Close()
	return err
}
