package repository

import (
	"time"

	"github.com/go-pg/pg/v10"
)

type Repo struct {
	db *pg.DB
}

func NewRepo(db *pg.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(task *TaskEntity) (*TaskEntity, error) {
	_, err := r.db.Model(task).Returning("*").Insert()
	return task, err
}

func (r *Repo) GetByID(id int64) (*TaskEntity, error) {
	task := &TaskEntity{}
	err := r.db.Model(task).Where("id = ?", id).First()
	return task, err
}

func (r *Repo) Update(id int64, task *TaskEntity) (*TaskEntity, error) {
	task.UpdatedAt = time.Now()
	_, err := r.db.Model(task).Where("id = ?", id).Returning("*").Update()
	if err != nil {
		return nil, err
	}
	return task, nil
}
