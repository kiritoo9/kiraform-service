package userrepo

import (
	"gorm.io/gorm"
)

type Query interface {
	FindByEmail(email string) (any, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) Query {
	return &Repository{DB: db}
}

func (r *Repository) FindByEmail(email string) (any, error) {
	return nil, nil
}
