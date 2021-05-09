package catalog

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"shopping/rest/content"
	"shopping/service"
)

type Catalog struct {
	Id   int    `form:"id"`
	Name string `form:"name"`
	Hide bool   `form:"hide"`
}

func (d *Catalog) serviceType() *service.Catalog {
	return &service.Catalog{
		Id:   d.Id,
		Name: d.Name,
		Hide: d.Hide,
	}
}

func (d *Catalog) Assemble(s *service.Catalog) {
	d.Id = s.Id
	d.Name = s.Name
	d.Hide = s.Hide
}

func FindCatalog(c *gin.Context) {
	page := &service.Page{}
	c.ShouldBind(page)
	catalog := &Catalog{}
	c.ShouldBind(catalog)

	rows, code := service.CatalogService.Find(catalog.serviceType(), page)
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

func RegisterCatalog(c *gin.Context) {
	catalog := &Catalog{}
	c.ShouldBind(catalog)
	m, code := service.CatalogService.Register(catalog.serviceType())
	resp := content.NewContent()
	resp.Data(m)
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func DeleteCatalog(c *gin.Context) {
	d, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "msg": "Page not found"})
		return
	}

	code := service.CatalogService.Delete(d)
	resp := content.NewContent()
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}

func UpdateCatalog(c *gin.Context) {
	catalog := &Catalog{}
	c.ShouldBind(catalog)
	affect, code := service.CatalogService.Update(catalog.serviceType())
	resp := content.NewContent()
	m := map[string]int64{
		"affect": affect,
	}
	resp.Data(m)
	resp.Code(code.Int()).Msg(service.ErrorMsg(code))
	c.JSON(service.ErrorStatusCode(code), resp)
}
