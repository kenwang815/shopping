package commodity

import (
	"time"

	"github.com/jinzhu/gorm"

	"shopping/model"
	"shopping/model/catalog"
)

type Commodity struct {
	Id          int             `gorm:"primaryKey;uniqueIndex;autoIncrement" mapKey:"ignore"`
	Catalog     catalog.Catalog `gorm:"ForeignKey:CatalogId;association_foreignkey:Id`
	CatalogId   int             `gorm:"column:catalog_id;not null;type:int" mapKey:"catalog_id,omitempty"`
	Name        string          `gorm:"column:name;type:varchar(64);not null" mapKey:"name,omitempty"`
	Cost        int             `gorm:"column:cost;type:int;not null" mapKey:"cost,omitempty"`
	Price       int             `gorm:"column:price;type:int;not null" mapKey:"price,omitempty"`
	Description string          `gorm:"column:description;type:text;not null" mapKey:"description,omitempty"`
	Sell        bool            `gorm:"column:sell;type:tinyint;not null" mapKey:"sell,omitempty"`
	StartTime   time.Time       `gorm:"column:start_time;type:datetime;not null" mapKey:"start_time,omitempty"`
	EndTime     time.Time       `gorm:"column:end_time;type:datetime;not null" mapKey:"end_time,omitempty"`
}

func (Commodity) TableName() string {
	return "commodity"
}

type Repository interface {
	Create(d *Commodity) (*Commodity, error)
	Update(d *Commodity) (*Commodity, int64, error)
	Delete(id int) (int64, error)
	Find(d *Commodity, p *model.Page) ([]*Commodity, error)
	Query(query interface{}, args ...interface{}) *gorm.DB
}
