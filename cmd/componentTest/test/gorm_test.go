package test

import (
	"fmt"
	"log"
	"mall/cmd/Account/entity"
	"mall/pkg/util"
	"testing"
)

const (
	gDSN = "root:123456@tcp(localhost:3306)/mall?charset=utf8&parseTime=True&loc=Local"
)

// 创建数据库表
func TestCreateTable(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	// 执行的建表语句是：CREATE TABLE `user_accounts` (`account_name` varchar(191),`passwordColumn` longtext,PRIMARY KEY (`account_name`))
	db.AutoMigrate(&entity.TbUserAccount{}) // 默认创建的表名是: user_accounts
}

// 往数据库表中插入一条记录，账号是root
func TestCreateTableRecordRoot(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	db.Create(&entity.TbUserAccount{
		AccountName: "root",
		Password:    "123456",
	}) // AccountName如果不是主键的话多次执行，会生成多条记录；否则只会执行一次，后续的执行会报错
}

// 往数据库表中插入一条记录，账号是leeBai
func TestCreateTableRecordLeeBai(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	db.Create(&entity.TbUserAccount{
		AccountName: "leeBai",
		Password:    "123456",
	}) // AccountName如果不是主键的话多次执行，会生成多条记录；否则只会执行一次，后续的执行会报错
}

// 读取数据库表中账号是root的记录
func TestReadTableRecordRoot(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	// be careful: 当数据库表中没有对应记录时,Find方法的Error返回是nil
	var userAccount entity.TbUserAccount
	err := db.Find(&userAccount, fmt.Sprintf("%s=?", entity.AccountNameColumn), "root").Error
	log.Printf("err:%v", err)
	log.Printf("query userAccount: %v\n", userAccount)
}

// 更新数据库表中账号是root的记录的密码
func TestUpdateTableRecordRoot(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	var userAccount entity.TbUserAccount
	// Where调用必须放在Update之前，否则无法更新
	db.Model(&userAccount).Where(fmt.Sprintf("%s=?", entity.AccountNameColumn), "root").Update(entity.PasswordColumn, "666666")
}

// 删除数据库表中账号是root的记录
func TestDeleteTableRecordRoot(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	var userAccount entity.TbUserAccount
	db.Delete(&userAccount, fmt.Sprintf("%s=?", entity.AccountNameColumn), "root")
}

// 删除数据库表中账号是leeBai的记录
func TestDeleteTableRecordLeeBai(t *testing.T) {
	db := util.InitMySQLDB(gDSN)

	var userAccount entity.TbUserAccount
	db.Delete(&userAccount, fmt.Sprintf("%s=?", entity.AccountNameColumn), "leeBai")
}
