package catalog

import (
	"shopping/model"

	"github.com/jinzhu/gorm"
)

type UUID string

func (u UUID) String() string { return string(u) }

type Catalog struct {
	Id   UUID   `gorm:"column:id;unique;type:uuid;primary_key" mapKey:"ignore"`
	Name string `gorm:"column:name;type:varchar(64);not null" mapKey:"name,omitempty"`
	Hide bool   `gorm:"column:hide;type:tinyint;not null" mapKey:"hide,omitempty"`
}

func (Catalog) TableName() string {
	return "catalog"
}

type Repository interface {
	Create(d *Catalog) (*Catalog, error)
	Update(d *Catalog) (*Catalog, int64, error)
	Delete(id UUID) (int64, error)
	Find(d *Catalog, p *model.Page) ([]*Catalog, error)
	Query(query interface{}, args ...interface{}) *gorm.DB
}
