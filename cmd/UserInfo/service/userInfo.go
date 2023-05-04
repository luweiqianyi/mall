package service

import (
	"mall/cmd/UserInfo/dao"
	"mall/cmd/UserInfo/entity"
)

func CreateUserInfo(info entity.TbUserInfo) error {
	return dao.CreateUserInfo(info)
}

func UpdateNickName(accountName string, nickName string) error {
	return dao.UpdateNickName(accountName, nickName)
}

func UpdatePortraitURL(accountName string, portrait string) error {
	return dao.UpdatePortraitURL(accountName, portrait)
}

func UpdateBirthday(accountName string, birthday string) error {
	return dao.UpdateBirthday(accountName, birthday)
}

func UpdatePhone(accountName string, phone string) error {
	return dao.UpdatePhone(accountName, phone)
}

func UpdateGender(accountName string, gender string) error {
	return dao.UpdateGender(accountName, gender)
}

func DeleteUserInfo(accountName string) error {
	return dao.DeleteUserInfo(accountName)
}

func QueryUserInfo(accountName string) (entity.TbUserInfo, error) {
	return dao.QueryUserInfo(accountName)
}
