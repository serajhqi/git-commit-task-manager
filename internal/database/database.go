package database

import (
	"git-project-management/config"
	"git-project-management/internal/types"
	"log"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var (
	db   *pg.DB
	once sync.Once
)

func GetDB() *pg.DB {

	config := config.GetConfig()
	once.Do(func() {
		db = pg.Connect(&pg.Options{
			Addr:     config.PG_HOST,
			User:     config.PG_USER,
			Password: config.PG_PASSWORD,
			Database: config.PG_DATABASE,
		})

		migrate(db)

	})

	return db
}

func migrate(db *pg.DB) error {
	models := []interface{}{
		(*types.ProjectEntity)(nil),
		(*types.TaskEntity)(nil),
		(*types.UserEntity)(nil),
		(*types.NotificationEntity)(nil),
		(*types.ActivityEntity)(nil),
		(*types.ApiKeyEntity)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true, // Skip if table already exists
		})
		if err != nil {
			log.Printf("Warning: could not create table for %T: %v", model, err)
		}
	}

	// Create indexes
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS idx_task_project_id ON tbl_task(project_id)",
		"CREATE INDEX IF NOT EXISTS idx_task_assignee_id ON tbl_task(assignee_id)",
		"CREATE INDEX IF NOT EXISTS idx_notification_user_id ON tbl_notification(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_activity_task_id ON tbl_activity(task_id)",
		"CREATE INDEX IF NOT EXISTS idx_activity_created_by ON tbl_activity(created_by)",
	}

	for _, query := range indexQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Warning: could not create index: %v", err)
		}
	}

	return nil
}
