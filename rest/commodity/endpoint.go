package commodity

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"shopping/rest/content"
	"shopping/service"
)

type Commodity struct {
	Id          int       `form:"id"`
	CatalogId   int       `form:"catalog_id"`
	Name        string    `form:"name"`
	Cost        int       `form:"cost"`
	Price       int       `form:"price"`
	Description string    `form:"description"`
	Sell        bool      `form:"sell"`
	StartTime   time.Time `form:"start_time" time_format:"2006-01-02 15:04:05"`
	EndTime     time.Time `form:"end_time" time_format:"2006-01-02 15:04:05"`
}

func (d *Commodity) serviceType() *service.Commodity {
	return &service.Commodity{
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

func (d *Commodity) Assemble(s *service.Commodity) {
	d.Id = s.Id
	d.CatalogId = s.CatalogId
	d.Name = s.Name
	d.Cost = s.Cost
	d.Price = s.Price
	d.Description = s.Description
	d.Sell = s.Sell
	d.StartTime = s.StartTime
	d.EndTime = s.EndTime
}

func FindCommodity(c *gin.Context) {
	page := &service.Page{}
	c.ShouldBind(page)
	commodity := &Commodity{}
	c.ShouldBind(commodity)

	rows, code := service.CommodityService.Find(commodity.serviceType(), page)
	resp := content.NewContent()
	if code == service.ErrorCodeSuccess {
		data := make(map[string]interface{})
		if page != nil {
			data["page"] = page
		}

		data["datas"] = rows
		resp.Data(data)
	}

	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func RegisterCommodity(c *gin.Context) {
	commodity := &Commodity{}
	c.ShouldBind(commodity)
	m, code := service.CommodityService.Register(commodity.serviceType())
	resp := content.NewContent()
	resp.Data(m)
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func DeleteCommodity(c *gin.Context) {
	d, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "msg": "Page not found"})
		return
	}

	code := service.CommodityService.Delete(d)
	resp := content.NewContent()
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func UpdateCommodity(c *gin.Context) {
	commodity := &Commodity{}
	c.ShouldBind(commodity)
	affect, code := service.CommodityService.Update(commodity.serviceType())
	resp := content.NewContent()
	m := map[string]int64{
		"affect": affect,
	}
	resp.Data(m)
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}
