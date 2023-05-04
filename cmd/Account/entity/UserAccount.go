package entity

const (
	AccountNameColumn = "accountName"
	PasswordColumn    = "password"
)

type TbUserAccount struct {
	// gorm.Model
	AccountName string `gorm:"primarykey;column:accountName"` //自定义数据库表字段名称为accountName
	Password    string `gorm:"column:password"`
}

// TableName 自定义数据库表名
func (u *TbUserAccount) TableName() string {
	return "TbUserAccount"
}
