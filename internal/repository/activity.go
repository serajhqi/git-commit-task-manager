package repository

import (
	"github.com/go-pg/pg/v10"
)



func (r *Repo) create(activity *type) (*ActivityEntity, error) {

	_, err := r.db.Model(activity).Returning("*").Insert()
	return activity, err
}

func (r *Repo) getByID(id int64) (*ActivityEntity, error) {
	project := &ActivityEntity{}
	err := r.db.Model(project).Where("id = ?", id).First()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repo) findByUserIDAndTaskID(userID, taskID int64) (*ActivityEntity, error) {
	project := &ActivityEntity{}
	err := r.db.Model(project).Where("created_by = ? AND task_id = ?", userID, taskID).Order("created_at DESC").First()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (r *Repo) getAll(taskID int64, limit, offset int) ([]ActivityEntity, error) {
	var activities []ActivityEntity

	err := r.db.Model(&activities).Where("task_id = ?", taskID).Limit(limit).Offset(offset).Select()

	if err != nil {
		return activities, err
	}
	return activities, nil

}

func (r *Repo) findByUserID(userID int64) (*ActivityEntity, error) {
	project := &ActivityEntity{}
	err := r.db.Model(project).Where("created_by = ?", userID).Order("created_at DESC").First()
	if err != nil {
		return nil, err
	}
	return project, nil
}
