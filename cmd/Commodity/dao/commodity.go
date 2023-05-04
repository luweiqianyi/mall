package dao

import (
	"fmt"
	"mall/cmd/Commodity/entity"
)

func AddCommodity(info entity.TbCommodity) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}

	record := info
	err := db.Create(&record).Error
	if err != nil {
		return fmt.Errorf("insert TbCommodity fail,err=%v", err)
	}
	return nil
}

func DeleteCommodity(identifyID string) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}

	var commodity entity.TbCommodity
	err := db.Delete(&commodity, fmt.Sprintf("%s=?", entity.IdentifyIDColumn), identifyID).Error
	return err
}

func UpdateCommodityColumn(identifyID string, columnName string, columnValue interface{}) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}
	var commodity entity.TbCommodity
	tx := db.Model(&commodity).Where(fmt.Sprintf("%s=?", entity.IdentifyIDColumn), identifyID).Update(columnName, columnValue)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("UpdateCommodityColumn[%s] failed: account[%s] not exist", columnName, identifyID)
	}
	if tx.Error != nil {
		return fmt.Errorf("UpdateCommodityColumn[%s] failed, err=%v", columnName, tx.Error)
	}
	return nil
}

func QueryCommodity(identifyID string) (entity.TbCommodity, error) {
	db := GetDB()
	if db == nil {
		return entity.TbCommodity{}, fmt.Errorf("%s", "db connection not open!")
	}
	var info entity.TbCommodity
	err := db.First(&info, fmt.Sprintf("%s=?", entity.IdentifyIDColumn), identifyID).Error
	if err != nil {
		return entity.TbCommodity{}, fmt.Errorf("QueryUserInfo failed,err=%v", err)
	}
	return info, nil
}

func QueryCommodities() ([]entity.TbCommodity, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("%s", "db connection not open!")
	}
	var commodities []entity.TbCommodity
	err := db.Table("TbCommodity").Find(&commodities).Error
	if err != nil {
		return nil, err
	}
	if len(commodities) == 0 {
		return nil, fmt.Errorf("QueryCommodities no result")
	}
	return commodities, nil
}

func UpdateCommodityMainCategory(identifyID string, mainCategory string) error {
	return UpdateCommodityColumn(identifyID, entity.MainCategoryColumn, mainCategory)
}

func UpdateCommoditySubCategory(identifyID string, subCategory string) error {
	return UpdateCommodityColumn(identifyID, entity.SubCategoryColumn, subCategory)
}

func UpdateCommodityName(identifyID string, name string) error {
	return UpdateCommodityColumn(identifyID, entity.NameColumn, name)
}

func UpdateCommodityPrice(identifyID string, price float64) error {
	return UpdateCommodityColumn(identifyID, entity.PriceColumn, price)
}
