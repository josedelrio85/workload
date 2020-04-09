package workload

import (
	"errors"
	"net/http"
	"time"
)

// Project represents project entity
type Project struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name"`
	Type      int64      `sql:"-" json:"-"`
}

var global []Project

// Global returns helper global variable
func (p *Project) Global() []Project {
	return global
}

// restartGlobal restarts value of global variable
func restartGlobal() {
	global = nil
}

// TableName sets the default table name
func (Project) TableName() string {
	return "project"
}

// GetType returns type of the element
func (p *Project) GetType() int64 {
	return p.Type
}

// GetElements retrieves an array of projectst
func (p *Project) GetElements(input Input, st Storer) error {
	restartGlobal()
	var ids []uint
	for _, p := range input.Project {
		ids = append(ids, p.ID)
	}

	a := st.Instance().Debug()
	if len(ids) == 0 {
		a.Find(&global)
	} else {
		a.Where("id in (?)", ids).Find(&global)
	}
	if a.Error != nil {
		return a.Error
	}
	return nil
}

// CreateElement creates a project
func (p *Project) CreateElement(input Input, st Storer) error {
	restartGlobal()
	for _, project := range input.Project {
		db := st.Instance().Debug()
		if err := db.Save(&project).Error; err != nil {
			return err
			// continue
		}
		global = append(global, project)
	}

	return nil
}

// UpdateElement updates a project
func (p *Project) UpdateElement(input Input, st Storer) error {
	restartGlobal()

	for _, project := range input.Project {
		db := st.Instance().Debug()
		update := db.Model(&project).Where("id = ?", project.ID).Updates(project)
		rows, err := update.RowsAffected, update.Error

		if err != nil {
			return err
		}

		if rows == 0 {
			err := errors.New("0 rows affected")
			return err
		}
		global = append(global, project)
	}
	return nil
}

// DeleteElement blablabla
func (p *Project) DeleteElement(input Input, st Storer) error {
	restartGlobal()

	for _, project := range input.Project {
		db := st.Instance().Debug()
		delete := db.Model(&project).Where("id = ?", project.ID).Delete(Project{})
		rows, err := delete.RowsAffected, delete.Error

		if err != nil {
			return err
		}

		if rows == 0 {
			err := errors.New("0 rows affected")
			return err
		}
	}
	return nil
}

// Response blablabla
func (p *Project) Response(w http.ResponseWriter) {

	resp := Response{
		Writer:  w,
		Code:    http.StatusOK,
		Message: "ALALALA",
		Project: p.Global(),
	}
	resp.Ok(p.Type)
}
