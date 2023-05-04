# gorm_demo
## 环境准备
1. 数据库环境准备：`mysql:5.7`
* `env`配置文件
    ```env
    CONTAINER_BASE_NAME=mall
    
    TZ=Asia/Shanghai
    
    #
    # mysql
    #
    MYSQL_ROOT_PASSWORD=123456
    DB_NAME=mall
    MYSQL_PORT=3306
    ```
* `docker-compose.yml`文件
    ```yaml
    version: "3.7"
    
    services:
      mysql:
        image: mysql:5.7
        container_name: ${CONTAINER_BASE_NAME}_mysql
        volumes:
          - user_login_volume:/var/lib/mysql
        environment:
          MYSQL_ROOT_PASSWORD: "${MYSQL_ROOT_PASSWORD}"
          MYSQL_DATABASE: "${DB_NAME}"
          TZ: "${TZ}"
        ports:
          - "${MYSQL_PORT}:3306"
    
    volumes:
      user_login_volume:
    ```
* 在`docker-compose.yml`文件所在路径执行脚本`docker-compose -p mall up -d`创建`docker`容器
  > -p: 指定名字
2. 下载框架库
  * `go get -u gorm.io/gorm`
  * `go get -u gorm.io/driver/mysql`
3. `mall/cmd/login/entity`目录下新建`UserAccount.go`，添加代码如下
  ```go
  package entity
  
  import "gorm.io/gorm"
  
  type UserAccount struct {
    gorm.Model
    Username string
    Password string
  }
  ```
4. `mall/cmd/login/test`目录下创建测试文件`gorm_test.go`,编写测试代码
  ```go
  package test
  
  import (
      "gorm.io/driver/mysql"
      "gorm.io/gorm"
      "log"
      "mall/cmd/login/entity"
      "testing"
  )
  
  func TestGorm(t *testing.T) {
      dsn := "root:123456@tcp(localhost:3306)/mall?charset=utf8&parseTime=True&loc=Local"
      db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
      if err != nil {
          panic(err)
      }
  
      // 执行的建表语句是：CREATE TABLE `user_accounts` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`username` longtext,`password` longtext,PRIMARY KEY (`id`),INDEX `idx_user_accounts_deleted_at` (`deleted_at`))
      db.AutoMigrate(&entity.UserAccount{}) // 默认创建的表名是: user_accounts
      db.Create(&entity.UserAccount{
          Username: "root",
          Password: "123456",
      })
  
      var user entity.UserAccount
      db.Find(&user, "username=?", "root")
      log.Printf("query: %v\n", user)
  
      db.Model(&user).Update("password", "666666")
      db.Find(&user, "username=?", "root")
      log.Printf("update: %v\n", user)
  }
  ```
  代码执行结果如下：
  ```log
  === RUN   TestGorm
  
  2023/04/21 11:18:07 F:/HappyCoding/goprograms/mall/cmd/login/test/gorm_test.go:18 SLOW SQL >= 200ms
  [654.802ms] [rows:0] CREATE TABLE `user_accounts` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`username` longtext,`password` longtext,PRIMARY KEY (`id`),INDEX `idx_user_accounts_deleted_at` (`deleted_at`))
  2023/04/21 11:18:07 query: {{1 2023-04-21 11:18:07.09 +0800 CST 2023-04-21 11:18:07.09 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} root 123456}
  2023/04/21 11:18:07 update: {{1 2023-04-21 11:18:07.09 +0800 CST 2023-04-21 11:18:07.141 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} root 666666}
  --- PASS: TestGorm (0.79s)
  PASS
  
  Process finished with the exit code 0
  ```
  与此同时`docker`容器:`mall_mysql`中数据库`mall`中新增了一个表，表中有了代码所产生的数据。详细情况如下所示：
  ```log
  sh-4.2# mysql -u root -p
  Enter password: 
  Welcome to the MySQL monitor.  Commands end with ; or \g.
  Your MySQL connection id is 2
  Server version: 5.7.41 MySQL Community Server (GPL)
  
  Copyright (c) 2000, 2023, Oracle and/or its affiliates.
  
  Oracle is a registered trademark of Oracle Corporation and/or its
  affiliates. Other names may be trademarks of their respective
  owners.
  
  Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
  
  mysql> show databases;
  +--------------------+
  | Database           |
  +--------------------+
  | information_schema |
  | mall               |
  | mysql              |
  | performance_schema |
  | sys                |
  +--------------------+
  5 rows in set (0.00 sec)
  
  mysql> use mall
  Reading table information for completion of table and column names
  You can turn off this feature to get a quicker startup with -A
  
  Database changed
  mysql> show tables;
  +----------------+
  | Tables_in_mall |
  +----------------+
  | user_accounts  |
  +----------------+
  1 row in set (0.00 sec)
  
  mysql> ^C
  mysql> select * from user_accounts;
  +----+-------------------------+-------------------------+------------+----------+----------+
  | id | created_at              | updated_at              | deleted_at | username | password |
  +----+-------------------------+-------------------------+------------+----------+----------+
  |  1 | 2023-04-21 11:18:07.090 | 2023-04-21 11:18:07.141 | NULL       | root     | 666666   |
  +----+-------------------------+-------------------------+------------+----------+----------+
  1 row in set (0.00 sec)
  ```
