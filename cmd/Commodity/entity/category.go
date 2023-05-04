package entity

const (
	MainCategory1 = "MC1"
	MainCategory2 = "MC2"
	MainCategory3 = "MC3"
	MainCategory4 = "MC4"

	SubCategory1 = "SC1"
	SubCategory2 = "SC2"
	SubCategory3 = "SC3"
)

type TbCategory struct {
	MainCategory string `gorm:"primaryKey;size:16;column:mainCategory"`
	SubCategory  string `gorm:"primaryKey;size:16;column:subCategory"`
} //指定复合主键

func (c *TbCategory) TableName() string {
	return "TbCategory"
}
