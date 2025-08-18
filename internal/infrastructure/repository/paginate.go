package repository

import (
	"math"

	"gorm.io/gorm"

	"itmrchow/go-todolist-service/internal/domain/repository"
)

func Paginate[T any](value interface{}, p *repository.Pagination[T], db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	// query count
	var totalRows int64
	db.Model(value).Count(&totalRows)
	p.TotalRows = totalRows

	// total page
	totalPages := int(math.Ceil(float64(totalRows) / float64(p.Limit)))
	p.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetSort())
	}
}
