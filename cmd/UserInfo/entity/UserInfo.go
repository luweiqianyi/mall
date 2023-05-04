package entity

const (
	AccountNameColumn = "accountName"
	NickNameColumn    = "nickName"
	PortraitURLColumn = "portraitURL"
	BirthdayColumn    = "birthday"
	PhoneColumn       = "phone"
	GenderColumn      = "gender"
)

type TbUserInfo struct {
	ID          uint   `gorm:"primarykey;column:id"`
	AccountName string `gorm:"unique;ForeignKey:accountName;column:accountName"`

	NickName    string `gorm:"unique;column:nickName"`
	PortraitURL string `gorm:"column:portraitURL"`
	Birthday    string `gorm:"column:birthday"`
	Phone       string `gorm:"column:phone"`
	Gender      string `gorm:"column:gender"`
}

func (u *TbUserInfo) TableName() string {
	return "TbUserInfo"
}
