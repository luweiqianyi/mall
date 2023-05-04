package dao

import (
	"fmt"
	"mall/cmd/Account/entity"
	"mall/pkg/log"
)

func Register(accountName string, password string) bool {
	// TODO 对用户的密码进行加密处理（使用不可逆加密算法，防止被第三方反向解密从而获取到用户的明文密码），再存入数据库，
	// 防止某些居心叵测的数据库维护人员拿到用户的账号密码为所欲为
	db := GetDB()
	if db == nil {
		return false
	}

	err := db.Create(&entity.TbUserAccount{
		AccountName: accountName,
		Password:    password,
	}).Error
	if err != nil {
		log.PrintLog("Register: account[%s] register failed,err:%v\n", accountName, err)
		return false
	}
	return true
}

func UnRegister(accountName string) bool {
	db := GetDB()
	if db == nil {
		return false
	}

	var userAccount entity.TbUserAccount
	err := db.Delete(&userAccount, fmt.Sprintf("%s=?", entity.AccountNameColumn), accountName).Error
	if err != nil {
		log.PrintLog("UnRegister account[%s] unregister failed,err:%v\n", accountName, err)
		return false
	}
	return true
}

func ChangePassword(accountName string, password string) bool {
	db := GetDB()
	if db == nil {
		return false
	}

	var userAccount entity.TbUserAccount
	// Where调用必须放在Update之前，否则无法更新
	tx := db.Model(&userAccount).Where(fmt.Sprintf("%s=?", entity.AccountNameColumn), accountName).Update(entity.PasswordColumn, password)
	if tx.RowsAffected == 0 {
		log.PrintLog("ChangePassword: account[%s] not exist!!", accountName)
		return false
	}
	if tx.Error != nil {
		log.PrintLog("ChangePassword: account[%s] password update failed,err=%v", accountName, tx.Error)
		return false
	}
	return true
}

func Login(accountName string, password string) bool {
	// TODO 若Register接口对用户密码进行加密，则对用户输入的明文密码使用Register接口中相同的加密算法进行加密，得到密文，
	// 密文与数据库中的匹配，则认为登录成功
	db := GetDB()
	if db == nil {
		return false
	}
	var userAccount entity.TbUserAccount
	// 不要用Find, 数据库中没有记录时Find返回的err为空
	err := db.First(&userAccount, fmt.Sprintf("%s=? and %s=?", entity.AccountNameColumn, entity.PasswordColumn), accountName, password).Error
	if err != nil {
		log.PrintLog("Login: account[%s] login failed,err:%v\n", accountName, err)
		return false
	}
	return true
}
