package service

import (
	"time"

	"shopping/model"
	"shopping/model/commodity"
	"shopping/utils/log"
)

type Commodity struct {
	Id          int       `json:"id"`
	CatalogId   int       `json:"catalog_id"`
	Name        string    `json:"name"`
	Cost        int       `json:"cost"`
	Price       int       `json:"price"`
	Description string    `json:"description"`
	Sell        bool      `json:"sell"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

func (d *Commodity) repoType() *commodity.Commodity {
	return &commodity.Commodity{
		Id:          d.Id,
		CatalogId:   d.CatalogId,
		Name:        d.Name,
		Cost:        d.Cost,
		Price:       d.Price,
		Description: d.Description,
		Sell:        d.Sell,
		StartTime:   d.StartTime,
		EndTime:     d.EndTime,
	}
}

func (d *Commodity) Assemble(r *commodity.Commodity) {
	d.Id = r.Id
	d.CatalogId = r.CatalogId
	d.Name = r.Name
	d.Cost = r.Cost
	d.Price = r.Price
	d.Description = r.Description
	d.Sell = r.Sell
	d.StartTime = r.StartTime
	d.EndTime = r.EndTime
}

type commodityService struct {
	commodityRepo commodity.Repository
}

func (s *commodityService) Find(d *Commodity, page *Page) ([]*Commodity, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	p := &model.Page{}
	if page != nil {
		p.Limit = page.Number
		p.Offset = page.Number * (page.Page - 1)
	}

	rows, err := s.commodityRepo.Find(d.repoType(), p)
	if err != nil {
		return nil, ErrorCodeCommodityDBFindFail
	}

	if len(rows) == 0 {
		return nil, ErrorCodeSuccessButNotFound
	}

	commoditys := []*Commodity{}
	for _, v := range rows {
		srv := &Commodity{}
		srv.Assemble(v)
		commoditys = append(commoditys, srv)
	}

	return commoditys, ErrorCodeSuccess
}

func (s *commodityService) Register(d *Commodity) (*Commodity, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	defaultTime := time.Time{}
	if d.Name == "" || d.Description == "" || d.StartTime == defaultTime || d.EndTime == defaultTime {
		return nil, ErrorCodeCommodityDBCreateFail
	}

	x, err := s.commodityRepo.Create(d.repoType())
	if err != nil {
		return nil, ErrorCodeCommodityDBCreateFail
	}

	re := &Commodity{}
	re.Assemble(x)
	return re, ErrorCodeSuccess
}

func (s *commodityService) Update(d *Commodity) (int64, ErrorCode) {
	_, a, err := s.commodityRepo.Update(d.repoType())
	if err != nil {
		return 0, ErrorCodeCommodityDBUpdateFail
	}

	return a, ErrorCodeSuccess
}

func (s *commodityService) Delete(id int) ErrorCode {
	affect, err := s.commodityRepo.Delete(id)
	if err != nil {
		return ErrorCodeCommodityDBDeleteFail
	}

	if affect <= 0 {
		log.Error("commodity does not exist: ", id)
		return ErrorCodeSuccessButNotFound
	}

	return ErrorCodeSuccess
}

func NewCommodityService(cr commodity.Repository) ICommodityService {
	return &commodityService{
		commodityRepo: cr,
	}
}
