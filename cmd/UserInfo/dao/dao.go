package dao

import (
	"fmt"
	"mall/cmd/UserInfo/entity"
)

func CreateUserInfo(info entity.TbUserInfo) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}

	record := info
	err := db.Create(&record).Error
	if err != nil {
		return fmt.Errorf("insert TbUserInfo fail,err=%v", err)
	}
	return nil
}

func UpdateNickName(accountName string, nickName string) error {
	return UpdateUserInfoColumn(accountName, entity.NickNameColumn, nickName)
}

func UpdatePortraitURL(accountName string, portrait string) error {
	return UpdateUserInfoColumn(accountName, entity.PortraitURLColumn, portrait)
}

func UpdateBirthday(accountName string, birthday string) error {
	return UpdateUserInfoColumn(accountName, entity.BirthdayColumn, birthday)
}

func UpdatePhone(accountName string, phone string) error {
	return UpdateUserInfoColumn(accountName, entity.PhoneColumn, phone)
}

func UpdateGender(accountName string, gender string) error {
	return UpdateUserInfoColumn(accountName, entity.GenderColumn, gender)
}

func UpdateUserInfoColumn(accountName string, columnName string, columnValue interface{}) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}
	var userInfo entity.TbUserInfo
	tx := db.Model(&userInfo).Where(fmt.Sprintf("%s=?", entity.AccountNameColumn), accountName).Update(columnName, columnValue)
	if tx.RowsAffected == 0 {
		return fmt.Errorf("UpdateUserInfoColumn[%s] failed: account[%s] not exist", columnName, accountName)
	}
	if tx.Error != nil {
		return fmt.Errorf("UpdateUserInfoColumn[%s] failed, err=%v", columnName, tx.Error)
	}
	return nil
}

func DeleteUserInfo(accountName string) error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("%s", "db connection not open!")
	}
	var userInfo entity.TbUserInfo
	err := db.Delete(&userInfo, fmt.Sprintf("%s=?", entity.AccountNameColumn), accountName).Error
	if err != nil {
		return fmt.Errorf("DeleteUserInfo[%s] failed,err=%v", accountName, err)
	}
	return nil
}

func QueryUserInfo(accountName string) (entity.TbUserInfo, error) {
	db := GetDB()
	if db == nil {
		return entity.TbUserInfo{}, fmt.Errorf("%s", "db connection not open!")
	}
	var userInfo entity.TbUserInfo
	err := db.First(&userInfo, fmt.Sprintf("%s=?", entity.AccountNameColumn), accountName).Error
	if err != nil {
		return entity.TbUserInfo{}, fmt.Errorf("QueryUserInfo failed,err=%v", err)
	}
	return userInfo, nil
}
