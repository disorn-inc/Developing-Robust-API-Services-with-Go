package store

import (
	"github/disorn-inc/Developing-Robust-API-Services-with-Go/Todo/todo"

	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *todo.Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) Find(todos *[]todo.Todo) error {
	return s.db.Find(todos).Error
}

func (s *GormStore) Delete(t *todo.Todo, id int) error {
	return s.db.Delete(t, id).Error
}