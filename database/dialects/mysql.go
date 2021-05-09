package dialects

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"shopping/config"
	"shopping/utils/log"
)

func MySQL(c *config.Database) *gorm.DB {
	connect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", c.User, c.Password, c.Host, c.Port, c.Name)

	myDB, err := gorm.Open("mysql", connect)
	if err != nil {
		log.Errorf("failed to connect database, %+v", err)
		return nil
	}

	return myDB
}
