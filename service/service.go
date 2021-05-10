package service

import (
	"shopping/config"
	"shopping/daos"
	"shopping/model/catalog"
	"shopping/model/commodity"
	"shopping/repository"
	"shopping/utils/log"
)

var (
	// === Repository ===
	CatalogRepo   catalog.Repository
	CommodityRepo commodity.Repository

	// === Service ===
	CatalogService   ICatalogService
	CommodityService ICommodityService
)

func Init(cf *config.Config, engine *repository.Engine) error {
	// === Repository ===
	CatalogRepo = daos.NewCatalogRepo(engine.GormDB)
	CommodityRepo = daos.NewCommodityRepo(engine.GormDB)

	// === Service ===
	CatalogService = NewCatalogService(CatalogRepo, CommodityRepo)
	CommodityService = NewCommodityService(CommodityRepo)

	log.Info("Create service success")
	return nil
}
