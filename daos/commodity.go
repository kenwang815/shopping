package daos

import (
	"github.com/jinzhu/gorm"

	"shopping/model"
	"shopping/model/commodity"
	"shopping/utils"
	"shopping/utils/log"
)

type commodityRepo struct {
	copy *gorm.DB
	db   *gorm.DB
}

func (r *commodityRepo) Create(d *commodity.Commodity) (*commodity.Commodity, error) {
	md := r.db.Select("catalog_id", "name", "cost", "price", "description", "sell", "start_time", "end_time").Create(d)
	if err := md.Error; err != nil {
		log.Errorf("commodityRepository Create fail => %+v", err)
		return nil, err
	}
	x := md.Value.(*commodity.Commodity)
	return x, nil
}

func (r *commodityRepo) Update(d *commodity.Commodity) (*commodity.Commodity, int64, error) {
	if d == nil {
		return nil, 0, nil
	}

	c := commodity.Commodity{
		Id: d.Id,
	}

	if *d == c {
		return nil, 0, nil
	}

	umap := utils.Map(d)

	re := &commodity.Commodity{}
	x := r.Query("id = ?", d.Id).Model(re).Updates(umap)
	if err := x.Error; err != nil {
		log.Errorf("commodityRepo update Info error: %+v", err)
		return nil, 0, err
	}
	affectRow := x.RowsAffected

	if affectRow <= 0 {
		return re, affectRow, x.Error
	}

	return re, affectRow, nil
}

func (r *commodityRepo) Delete(id int) (int64, error) {
	delete := r.db.Where("id = ?", id).Delete(&commodity.Commodity{})
	row := delete.RowsAffected
	var err error
	if err = delete.Error; err != nil {
		log.Errorf("commodityRepo Delete fail => %+v", err)
	}
	return row, err
}

func (r *commodityRepo) Find(d *commodity.Commodity, p *model.Page) ([]*commodity.Commodity, error) {
	var catalogs []*commodity.Commodity
	if err := r.db.Where(d).Limit(p.Limit).Offset(p.Offset).Find(&catalogs).Error; err != nil {
		log.Errorf("commodityRepo Find fail => %+v", err)
		return nil, err
	}
	return catalogs, nil
}

func (r *commodityRepo) Query(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func (r *commodityRepo) NewTransactions() {
	r.db = r.db.Begin()
}

func (r *commodityRepo) TransactionsRollback() {
	r.db.Rollback()
	r.db = r.copy
}

func (r *commodityRepo) TransactionsCommit() {
	r.db.Commit()
	r.db = r.copy
}

func NewCommodityRepo(db *gorm.DB) commodity.Repository {
	return &commodityRepo{
		copy: db,
		db:   db,
	}
}
