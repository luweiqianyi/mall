package entity

const (
	IdentifyIDColumn     = "identifyID"
	MainCategoryColumn   = "mainCategory"
	SubCategoryColumn    = "subCategory"
	NameColumn           = "name"
	PriceColumn          = "price"
	ProductionDateColumn = "productionDate"
	ExpireDateColumn     = "expireDate"
)

type CommodityProperty struct {
	Name  string  `gorm:"size:255;column:name"`
	Price float64 `gorm:"column:price"`
	PD    string  `gorm:"size:16;column:productionDate"`
	Exp   string  `gorm:"size:8;column:expireDate"`
}

type TbCommodity struct {
	IdentifyID   string `gorm:"primarykey;size:128;column:identifyID"` // 商品识别码，唯一
	MainCategory string `gorm:"size:16;column:mainCategory"`
	SubCategory  string `gorm:"size:16;column:subCategory"`
	CommodityProperty
}

func (c *TbCommodity) TableName() string {
	return "TbCommodity"
}
