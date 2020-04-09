package workload

import (
	"errors"
	"net/http"
	"time"
)

// Work represents work entity
type Work struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`

	PersonID uint   `json:"person_id"`
	Person   Person `gorm:"ForeignKey:PersonID;AssociationForeignKey:ID"`

	ProjectID uint    `json:"project_id"`
	Project   Project `gorm:"ForeignKey:ProjectID;AssociationForeignKey:ID"`

	StatusID uint   `json:"status_id"`
	Status   Status `gorm:"ForeignKey:StatusID;AssociationForeignKey:ID"`

	Description   string  `json:"description,omitempty"`
	EstimatedTime float32 `json:"estimated_time,omitempty"`
	RealTime      float32 `json:"real_time,omitempty"`
	Type          int64   `sql:"-" json:"-"`
}

// TableName sets the default table name
func (Work) TableName() string {
	return "work"
}

var globalwork []Work

// Global returns helper global variable
func (p *Work) Global() []Work {
	return globalwork
}

// restartWorkGlobal restarts value of global variable
func restartWorkGlobal() {
	globalwork = nil
}

// GetType returns type of the element
func (p *Work) GetType() int64 {
	return p.Type
}

// GetElements retrieves an array of works
func (p *Work) GetElements(input Input, st Storer) error {
	restartWorkGlobal()
	var ids []uint
	for _, p := range input.Project {
		ids = append(ids, p.ID)
	}

	a := st.Instance().Debug()
	if len(ids) == 0 {
		a.Preload("Project").Preload("Status").Preload("Person").Find(&globalwork)
	} else {
		a.Where("id in (?)", ids).Preload("Project").Preload("Status").Preload("Person").Find(&globalwork)
	}

	if a.Error != nil {
		return a.Error
	}
	return nil
}

// CreateElement creates a work
func (p *Work) CreateElement(input Input, st Storer) error {
	restartWorkGlobal()
	for _, work := range input.Work {
		db := st.Instance().Debug()
		if err := db.Save(&work).Error; err != nil {
			return err
		}
		globalwork = append(globalwork, work)
	}
	return nil
}

// UpdateElement updates a work
func (p *Work) UpdateElement(input Input, st Storer) error {
	restartWorkGlobal()

	for _, work := range input.Work {
		db := st.Instance().Debug()
		update := db.Model(&work).Where("id = ?", work.ID).Updates(work)
		rows, err := update.RowsAffected, update.Error

		if err != nil {
			return err
		}

		if rows == 0 {
			err := errors.New("0 rows affected")
			return err
		}
		globalwork = append(globalwork, work)
	}
	return nil
}

// DeleteElement deletes a work
func (p *Work) DeleteElement(input Input, st Storer) error {
	restartWorkGlobal()

	for _, work := range input.Work {
		db := st.Instance().Debug()
		delete := db.Model(&work).Where("id = ?", work.ID).Delete(Work{})
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
func (p *Work) Response(w http.ResponseWriter) {

	resp := Response{
		Writer: w,
		Code:   http.StatusOK,
		// Message: "ALALALA",
		Work: p.Global(),
	}
	resp.Ok(p.Type)
}
