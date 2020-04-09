package workload

import (
	"errors"
	"log"
	"net/http"
	"time"
)

// Person represents person entity
type Person struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name,omitempty"`
	Type      int64      `sql:"-" json:"-"`
}

// TableName sets the default table name
func (Person) TableName() string {
	return "person"
}

var globalperson []Person

// Global returns helper global variable
func (p *Person) Global() []Person {
	return globalperson
}

// restartstatusGlobal restarts value of global variable
func restartPersonGlobal() {
	globalperson = nil
}

// GetType returns type of the element
func (p *Person) GetType() int64 {
	return p.Type
}

// GetElements retrieves an array of person
func (p *Person) GetElements(input Input, st Storer) error {
	restartGlobal()
	var ids []uint
	for _, p := range input.Status {
		ids = append(ids, p.ID)
	}

	a := st.Instance().Debug()
	if len(ids) == 0 {
		a.Find(&globalperson)
	} else {
		a.Where("id in (?)", ids).Find(&globalperson)
	}
	if a.Error != nil {
		return a.Error
	}
	return nil
}

// CreateElement creates a person
func (p *Person) CreateElement(input Input, st Storer) error {
	restartGlobal()
	for _, person := range input.Person {
		log.Println(person)
		db := st.Instance().Debug()
		if err := db.Save(&person).Error; err != nil {
			return err
		}
		globalperson = append(globalperson, person)
	}
	return nil
}

// UpdateElement updates a person
func (p *Person) UpdateElement(input Input, st Storer) error {
	restartGlobal()

	for _, person := range input.Person {
		db := st.Instance().Debug()
		update := db.Model(&person).Where("id = ?", person.ID).Updates(person)
		rows, err := update.RowsAffected, update.Error

		if err != nil {
			return err
		}

		if rows == 0 {
			err := errors.New("0 rows affected")
			return err
		}
		globalperson = append(globalperson, person)
	}
	return nil
}

// DeleteElement deletes a person
func (p *Person) DeleteElement(input Input, st Storer) error {
	restartGlobal()

	for _, person := range input.Person {
		db := st.Instance().Debug()
		delete := db.Model(&person).Where("id = ?", person.ID).Delete(Person{})
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
func (p *Person) Response(w http.ResponseWriter) {

	resp := Response{
		Writer:  w,
		Code:    http.StatusOK,
		Message: "ALALALA",
		Person:  p.Global(),
	}
	resp.Ok(p.Type)
}
