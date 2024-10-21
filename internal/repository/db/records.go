package db

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cewuandy/go-restful-dns/internal/domain"
	"github.com/cewuandy/go-restful-dns/internal/repository/db/models"
	"github.com/cewuandy/go-restful-dns/internal/utils"

	"github.com/pkg/errors"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type recordRepo struct {
	db *gorm.DB
}

func (r *recordRepo) Create(ctx context.Context, record *domain.Record) error {
	var (
		raw models.Record
		err error
	)

	_ = utils.Convert(&record, &raw)
	err = r.db.WithContext(ctx).Create(&raw).Error
	if err != nil {
		return &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(err.Error()),
		}
	}

	return nil
}

func (r *recordRepo) Get(ctx context.Context, name string, rrType uint16,
	class uint16) (*domain.Record, error) {
	var (
		raw    models.Record
		record domain.Record
		err    error
	)

	err = r.db.WithContext(ctx).
		Unscoped().
		Where("name=? AND rr_type=? AND class=?", name, rrType, class).
		First(&raw).
		Error
	if err != nil {
		return nil, &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(err.Error()),
		}
	}

	_ = utils.Convert(&raw, &record)

	return &record, nil
}

func (r *recordRepo) List(ctx context.Context) ([]*domain.Record, error) {
	var (
		raws    []models.Record
		records []*domain.Record
		err     error
	)

	err = r.db.WithContext(ctx).Find(&raws).Error
	if err != nil {
		return nil, err
	}

	for _, raw := range raws {
		record := domain.Record{}
		_ = utils.Convert(&raw, &record)
		records = append(records, &record)
	}

	return records, nil
}

func (r *recordRepo) Update(ctx context.Context, record *domain.Record) error {
	var (
		raw models.Record
		err error
	)

	err = r.db.
		WithContext(ctx).
		Where("name=? AND rr_type=? AND class=?", record.Name, record.RrType, record.Class).
		First(&raw).
		Error
	if err != nil {
		return &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusNotFound,
			Err:        errors.New(err.Error()),
		}
	}

	_ = utils.Convert(&record, &raw)
	err = r.db.WithContext(ctx).Updates(&raw).Error
	if err != nil {
		return &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(err.Error()),
		}
	}

	return nil
}

func (r *recordRepo) Delete(ctx context.Context, name string, rrType uint16, class uint16) error {
	var (
		raw models.Record
		err error
	)

	err = r.db.
		WithContext(ctx).
		Where("name=? AND rr_type=? AND class=?", name, rrType, class).
		First(&raw).
		Error
	if err != nil {
		return &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusNotFound,
			Err:        errors.New(err.Error()),
		}
	}

	err = r.db.WithContext(ctx).
		Unscoped().
		Where("name=? AND rr_type=? AND class=?", name, rrType, class).
		Delete(&raw).
		Error
	if err != nil {
		return &domain.Error{
			Message:    fmt.Sprintf("DB error: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Err:        errors.New(err.Error()),
		}
	}

	return nil
}

func NewRecordsRepo(injector *do.Injector) (domain.RecordRepo, error) {
	return &recordRepo{do.MustInvoke[*gorm.DB](injector)}, nil
}
