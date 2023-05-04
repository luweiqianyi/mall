package dao

import (
	"fmt"
	"gorm.io/gorm"
	"mall/cmd/Account/config"
	"mall/cmd/Account/entity"
)

func GetDB() *gorm.DB {
	if config.GetGConfig() == nil || config.GetGConfig().GetDB() == nil {
		return nil
	}
	return config.GetGConfig().GetDB()
}

func CreateTables() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}

	// AutoMigrate：表结构数据类型发生变化，会自动修改相应字段
	db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&entity.TbUserAccount{})

	return nil
}
