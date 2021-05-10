package service

import (
	"time"

	"shopping/model"
	"shopping/model/catalog"
	"shopping/model/commodity"
	"shopping/utils/log"
)

type Catalog struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Hide bool   `json:"hide"`
}

func (d *Catalog) repoType() *catalog.Catalog {
	return &catalog.Catalog{
		Id:   d.Id,
		Name: d.Name,
		Hide: d.Hide,
	}
}

func (d *Catalog) Assemble(r *catalog.Catalog) {
	d.Id = r.Id
	d.Name = r.Name
	d.Hide = r.Hide
}

type catalogService struct {
	catalogRepo   catalog.Repository
	commodityRepo commodity.Repository
}

func (s *catalogService) Find(d *Catalog, page *Page) ([]*Catalog, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	p := &model.Page{}
	if page != nil {
		p.Limit = page.Number
		p.Offset = page.Number * (page.Page - 1)
	}

	rows, err := s.catalogRepo.Find(d.repoType(), p)
	if err != nil {
		return nil, ErrorCodeCatalogDBFindFail
	}

	if len(rows) == 0 {
		return nil, ErrorCodeSuccessButNotFound
	}

	catalogs := []*Catalog{}
	for _, v := range rows {
		srv := &Catalog{}
		srv.Assemble(v)
		catalogs = append(catalogs, srv)
	}

	return catalogs, ErrorCodeSuccess
}

func (s *catalogService) FindCommodity(i int, page *Page) ([]*Commodity, ErrorCode) {
	p := &model.Page{}
	if page != nil {
		p.Limit = page.Number
		p.Offset = page.Number * (page.Page - 1)
	}

	tmp, err := s.catalogRepo.Find(&catalog.Catalog{Id: i}, &model.Page{Limit: 1, Offset: 0})
	if err != nil {
		return nil, ErrorCodeCatalogDBFindFail
	}

	if len(tmp) == 0 {
		return nil, ErrorCodeSuccessButNotFound
	}

	if tmp[0].Hide {
		return nil, ErrorCodeNotFound
	}

	rows, err := s.commodityRepo.Find(&commodity.Commodity{CatalogId: i}, p)
	if err != nil {
		return nil, ErrorCodeCommodityDBFindFail
	}

	if len(rows) == 0 {
		return nil, ErrorCodeSuccessButNotFound
	}

	currentTime := time.Now().UTC()
	commoditys := []*Commodity{}
	for _, v := range rows {
		if v.Sell && v.StartTime.Before(currentTime) && v.EndTime.After(currentTime) {
			srv := &Commodity{}
			srv.Assemble(v)
			commoditys = append(commoditys, srv)
		}
	}

	return commoditys, ErrorCodeSuccess
}

func (s *catalogService) Register(d *Catalog) (*Catalog, ErrorCode) {
	if d == nil {
		return nil, ErrorCodeBadRequest
	}

	if d.Name == "" {
		return nil, ErrorCodeDataVerificationFail
	}

	x, err := s.catalogRepo.Create(d.repoType())
	if err != nil {
		return nil, ErrorCodeCatalogDBCreateFail
	}

	re := &Catalog{}
	re.Assemble(x)
	return re, ErrorCodeSuccess
}

func (s *catalogService) Update(d *Catalog) (int64, ErrorCode) {
	_, a, err := s.catalogRepo.Update(d.repoType())
	if err != nil {
		return 0, ErrorCodeCatalogDBUpdateFail
	}

	return a, ErrorCodeSuccess
}

func (s *catalogService) Delete(id int) ErrorCode {
	affect, err := s.catalogRepo.Delete(id)
	if err != nil {
		return ErrorCodeCatalogDBDeleteFail
	}
	if affect <= 0 {
		log.Error("catalog does not exist: ", id)
		return ErrorCodeSuccessButNotFound
	}

	return ErrorCodeSuccess
}

func NewCatalogService(cr catalog.Repository, cor commodity.Repository) ICatalogService {
	return &catalogService{
		catalogRepo:   cr,
		commodityRepo: cor,
	}
}
