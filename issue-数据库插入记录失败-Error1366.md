* 写入数据库的内容为中文时提示`Error 1366 (HY000): Incorrect string value: '\xE6\x9D\x8E\xE7\x99\xBD' for column 'nickName' at row 1`
解决方案: 创建数据库表时指定表为utf8,gorm的设置是在AutoMigrate调用前指定表的字符集：
```go
db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&entity.UserInfo{})
db.Create(&record)
```
检查是否生效，进入docker中的mysql Terminal，输入`show create table TbUserInfo;`查看建表语句,如下：
```log
mysql> show create table TbUserInfo;
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Table      | Create Table                                                                                                                                                                                                                                                                                                                                                                                                              |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| TbUserInfo | CREATE TABLE `TbUserInfo` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `accountName` varchar(191) DEFAULT NULL,
  `nickName` varchar(191) DEFAULT NULL,
  `portraitURL` longtext,
  `birthday` longtext,
  `phone` longtext,
  `gender` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `accountName` (`accountName`),
  UNIQUE KEY `nickName` (`nickName`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.00 sec)
```
说明go代码指定建表的字符为utf8确实生效了。如果建表时没有指定utf8,建表信息就是如下所示：
```log
```log
mysql> show create table TbUserInfo;
+------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| Table      | Create Table                                                                                                                                                                                                                                                                                                                                                                                               |
+------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
| TbUserInfo | CREATE TABLE `TbUserInfo` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `accountName` varchar(191) DEFAULT NULL,
  `nickName` varchar(191) DEFAULT NULL,
  `portraitURL` longtext,
  `birthday` longtext,
  `phone` longtext,
  `gender` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `accountName` (`accountName`),
  UNIQUE KEY `nickName` (`nickName`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 |
+------------+------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+
1 row in set (0.02 sec)
```
经过上述go代码修改，写入数据库成功，通过我的查询代码，查出来我的插入数据为
```log
{
    "code": 2013,
    "info": "{\"ID\":1,\"AccountName\":\"leebai\",\"NickName\":\"李白\",\"PortraitURL\":\"localhost:8080/path\",\"Birthday\":\"1990-01-01\",\"Phone\":\"13688449696\",\"Gender\":\"male\"}"
}
```
说明确实成功了。


对于业务流程来说确实已经解决了，没有问题了，此时我想要到我的docker的mysql容器中去看看我的插入数据是什么样的，查询结果如下：
```log
sh-4.2# mysql -u root -p 
Enter password: 
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 6
Server version: 5.7.41 MySQL Community Server (GPL)

Copyright (c) 2000, 2023, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql> use mall;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+----------------+
| Tables_in_mall |
+----------------+
| TbUserInfo     |
+----------------+
1 row in set (0.00 sec)

mysql> select * from TbUserInfo;
+----+-------------+----------+---------------------+------------+-------------+--------+
| id | accountName | nickName | portraitURL         | birthday   | phone       | gender |
+----+-------------+----------+---------------------+------------+-------------+--------+
|  1 | leebai      | ??       | localhost:8080/path | 1990-01-01 | 13688449696 | male   |
+----+-------------+----------+---------------------+------------+-------------+--------+
1 row in set (0.00 sec)
```
可以看到`nickName`所在列的中文变成了`??`，但是我们的go程序是可以返回正常的中文`李白`的，所以怀疑是mysql容器中的客户端连接工具的字符集有问题，
输入如下命令，查询客户端的字符集是`latin1`，该字符集无法显示中文字符。
```log
mysql> show variables like 'character%';
+--------------------------+----------------------------+
| Variable_name            | Value                      |
+--------------------------+----------------------------+
| character_set_client     | latin1                     |
| character_set_connection | latin1                     |
| character_set_database   | latin1                     |
| character_set_filesystem | binary                     |
| character_set_results    | latin1                     |
| character_set_server     | utf8mb4                    |
| character_set_system     | utf8                       |
| character_sets_dir       | /usr/share/mysql/charsets/ |
+--------------------------+----------------------------+
8 rows in set (0.00 sec)

mysql> show variables like 'collation%';
+----------------------+--------------------+
| Variable_name        | Value              |
+----------------------+--------------------+
| collation_connection | latin1_swedish_ci  |
| collation_database   | latin1_swedish_ci  |
| collation_server     | utf8mb4_unicode_ci |
+----------------------+--------------------+
3 rows in set (0.00 sec)
```
经过一番面向Google查询，只需要修改mysql容器的`/etc`目录下的`my.cnf`文件即可。具体修改步骤为在`[client]`项下加入`default-character-set=utf8`即可。修改后整个文件如下所示：
```log
# For advice on how to change settings please see
# http://dev.mysql.com/doc/refman/5.7/en/server-configuration-defaults.html

[mysqld]
#
# Remove leading # and set to the amount of RAM for the most important data
# cache in MySQL. Start at 70% of total RAM for dedicated server, else 10%.
# innodb_buffer_pool_size = 128M
#
# Remove leading # to turn on a very important data integrity option: logging
# changes to the binary log between backups.
# log_bin
#
# Remove leading # to set options mainly useful for reporting servers.
# The server defaults are faster for transactions and fast SELECTs.
# Adjust sizes as needed, experiment to find the optimal values.
# join_buffer_size = 128M
# sort_buffer_size = 2M
# read_rnd_buffer_size = 2M
skip-host-cache
skip-name-resolve
datadir=/var/lib/mysql
socket=/var/run/mysqld/mysqld.sock
secure-file-priv=/var/lib/mysql-files
user=mysql

# Disabling symbolic-links is recommended to prevent assorted security risks
symbolic-links=0

#log-error=/var/log/mysqld.log
pid-file=/var/run/mysqld/mysqld.pid
[client]
socket=/var/run/mysqld/mysqld.sock
default-character-set=utf8

!includedir /etc/mysql/conf.d/
!includedir /etc/mysql/mysql.conf.d/
```
修改之后，重启mysql容器，重新登录mysql，执行`select * from TbUserInfo;`,查询结果如下:
```log
mysql> select * from TbUserInfo;
+----+-------------+----------+---------------------+------------+-------------+--------+
| id | accountName | nickName | portraitURL         | birthday   | phone       | gender |
+----+-------------+----------+---------------------+------------+-------------+--------+
|  1 | leebai      | 李白     | localhost:8080/path | 1990-01-01 | 13688449696 | male   |
+----+-------------+----------+---------------------+------------+-------------+--------+
1 row in set (0.00 sec)
```
可以看到: `??`变成了`李白`，说明配置修改生效。

> 注意： 运行`docker-compose up -d`重新部署mysql容器会导致`my.cnf`中的配置项丢失。T_T