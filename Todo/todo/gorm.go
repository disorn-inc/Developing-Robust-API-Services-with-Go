package todo

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) Find(todos *[]Todo) error {
	return s.db.Find(todos).Error
}

func (s *GormStore) Delete(t *Todo, id int) error {
	return s.db.Delete(t, id).Error
}