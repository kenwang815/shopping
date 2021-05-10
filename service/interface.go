package service

type ICatalogService interface {
	Find(*Catalog, *Page) ([]*Catalog, ErrorCode)
	FindCommodity(int, *Page) ([]*Commodity, ErrorCode)
	Register(*Catalog) (*Catalog, ErrorCode)
	Update(*Catalog) (int64, ErrorCode)
	Delete(int) ErrorCode
}

type ICommodityService interface {
	Find(*Commodity, *Page) ([]*Commodity, ErrorCode)
	Register(*Commodity) (*Commodity, ErrorCode)
	Update(*Commodity) (int64, ErrorCode)
	Delete(int) ErrorCode
}
