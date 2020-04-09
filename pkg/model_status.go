package workload

import (
	"errors"
	"net/http"
	"time"
)

// Status represents status entity
type Status struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	Description string     `json:"description,omitempty"`
	Type        int64      `sql:"-" json:"-"`
}

// TableName sets the default table name
func (Status) TableName() string {
	return "status"
}

var globalstatus []Status

// Global returns helper global variable
func (p *Status) Global() []Status {
	return globalstatus
}

// restartstatusGlobal restarts value of global variable
func restartstatusGlobal() {
	globalstatus = nil
}

// GetType returns type of the element
func (p *Status) GetType() int64 {
	return p.Type
}

// GetElements retrieves an array of statuss
func (p *Status) GetElements(input Input, st Storer) error {
	restartGlobal()
	var ids []uint
	for _, p := range input.Status {
		ids = append(ids, p.ID)
	}

	a := st.Instance().Debug()
	if len(ids) == 0 {
		a.Find(&globalstatus)
	} else {
		a.Where("id in (?)", ids).Find(&globalstatus)
	}
	if a.Error != nil {
		return a.Error
	}
	return nil
}

// CreateElement creates a status
func (p *Status) CreateElement(input Input, st Storer) error {
	restartGlobal()
	for _, status := range input.Status {
		db := st.Instance().Debug()
		if err := db.Save(&status).Error; err != nil {
			return err
		}
		globalstatus = append(globalstatus, status)
	}
	return nil
}

// UpdateElement updates a status
func (p *Status) UpdateElement(input Input, st Storer) error {
	restartGlobal()

	for _, status := range input.Status {
		db := st.Instance().Debug()
		update := db.Model(&status).Where("id = ?", status.ID).Updates(status)
		rows, err := update.RowsAffected, update.Error

		if err != nil {
			return err
		}

		if rows == 0 {
			err := errors.New("0 rows affected")
			return err
		}
		globalstatus = append(globalstatus, status)
	}
	return nil
}

// DeleteElement deletes a status
func (p *Status) DeleteElement(input Input, st Storer) error {
	restartGlobal()

	for _, status := range input.Status {
		db := st.Instance().Debug()
		delete := db.Model(&status).Where("id = ?", status.ID).Delete(Status{})
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
func (p *Status) Response(w http.ResponseWriter) {

	resp := Response{
		Writer:  w,
		Code:    http.StatusOK,
		Message: "ALALALA",
		Status:  p.Global(),
	}
	resp.Ok(p.Type)
}
