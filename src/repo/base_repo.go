package repo

import "gorm.io/gorm"

type BaseRepo struct {
	repo *gorm.DB
}

type BaseRepoMethods interface {
	BeginTransaction(tx func(conn *gorm.DB) error) error
	SetRepo(repo *gorm.DB)
}

func (r *BaseRepo) SetRepo(repo *gorm.DB) {
	r.repo = repo
}

func (r *BaseRepo) BeginTransaction(tx func(conn *gorm.DB) error) error {
	return r.repo.Transaction(tx)
}
