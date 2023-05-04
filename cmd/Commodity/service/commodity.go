package service

import (
	"mall/cmd/Commodity/dao"
	"mall/cmd/Commodity/entity"
)

func AddCommodity(info entity.TbCommodity) error {
	return dao.AddCommodity(info)
}

func DeleteCommodity(identifyID string) error {
	return dao.DeleteCommodity(identifyID)
}

func QueryCommodity(identifyID string) (entity.TbCommodity, error) {
	return dao.QueryCommodity(identifyID)
}

func QueryCommodities() ([]entity.TbCommodity, error) {
	return dao.QueryCommodities()
}

func UpdateCommodityMainCategory(identifyID string, mainCategory string) error {
	return dao.UpdateCommodityMainCategory(identifyID, mainCategory)
}

func UpdateCommoditySubCategory(identifyID string, subCategory string) error {
	return dao.UpdateCommoditySubCategory(identifyID, subCategory)
}

func UpdateCommodityName(identifyID string, name string) error {
	return dao.UpdateCommodityName(identifyID, name)
}

func UpdateCommodityPrice(identifyID string, price float64) error {
	return dao.UpdateCommodityPrice(identifyID, price)
}
