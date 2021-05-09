package daos

import (
	"github.com/jinzhu/gorm"

	"shopping/model"
	"shopping/model/catalog"
	"shopping/utils"
	"shopping/utils/log"
)

type catalogRepo struct {
	copy *gorm.DB
	db   *gorm.DB
}

func (r *catalogRepo) Create(d *catalog.Catalog) (*catalog.Catalog, error) {
	md := r.db.Create(d)
	if err := md.Error; err != nil {
		log.Errorf("catalogRepository Create fail => %+v", err)
		return nil, err
	}
	x := md.Value.(*catalog.Catalog)
	return x, nil
}

func (r *catalogRepo) Update(d *catalog.Catalog) (*catalog.Catalog, int64, error) {
	if d == nil {
		return nil, 0, nil
	}

	c := catalog.Catalog{
		Id: d.Id,
	}

	if *d == c {
		return nil, 0, nil
	}

	umap := utils.Map(d)

	re := &catalog.Catalog{}
	x := r.Query("id = ?", d.Id).Model(re).Updates(umap)
	if err := x.Error; err != nil {
		log.Errorf("catalogRepository update Info error: %+v", err)
		return nil, 0, err
	}
	affectRow := x.RowsAffected

	if affectRow <= 0 {
		return re, affectRow, x.Error
	}

	return re, affectRow, nil
}

func (r *catalogRepo) Delete(id catalog.UUID) (int64, error) {
	delete := r.db.Where("id = ?", id).Delete(&catalog.Catalog{})
	row := delete.RowsAffected
	var err error
	if err = delete.Error; err != nil {
		log.Errorf("catalogRepository Delete fail => %+v", err)
	}
	return row, err
}

func (r *catalogRepo) Find(d *catalog.Catalog, p *model.Page) ([]*catalog.Catalog, error) {
	var catalogs []*catalog.Catalog
	if err := r.db.Where(d).Limit(p.Limit).Offset(p.Offset).Find(&catalogs).Error; err != nil {
		log.Errorf("catalogRepository Find fail => %+v", err)
		return nil, err
	}
	return catalogs, nil
}

func (r *catalogRepo) Query(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func (r *catalogRepo) NewTransactions() {
	r.db = r.db.Begin()
}

func (r *catalogRepo) TransactionsRollback() {
	r.db.Rollback()
	r.db = r.copy
}

func (r *catalogRepo) TransactionsCommit() {
	r.db.Commit()
	r.db = r.copy
}

func NewCatalogRepo(db *gorm.DB) catalog.Repository {
	return &catalogRepo{
		copy: db,
		db:   db,
	}
}
