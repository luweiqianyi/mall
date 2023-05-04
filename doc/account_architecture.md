# Account Architecture
与客户端的通信方式，支持http和https。
## 基本功能
用户账号，用户使用账号登录成功可以访问应用内数据。支持以下功能。
* 账号注册
* 修改账号密码
* 账号注销
* 用户登录
    * 允许同一账号在多端登录
    * 允许同一账号多端同时在线
* 用户退出
> 用户基数庞大时，以上操作对服务有性能压力 TODO 优化

## 路由设计
* URL访问地址支持http和https。
* 访问域名：`schema://domain`,比如`http://mycompany.com`或者`https://mycompany.com`
* 接口地址设计
    * `/api/user/register`: 账号注册
    * `/api/user/unregister`: 账号注销
    * `/api/user/login`: 用户登录
    * `/api/user/logout`: 用户退出

## 接口设计
### /api/user/register
账号注册

| 参数名       | 参数类型  | 参数说明|
|-----------| ----  |----|
| accountName  | string |账户名|
| password  | string |密码|

### /api/user/unregister
账号注销

| 参数名       | 参数类型  | 参数说明|
|-----------| ----  |----|
| accountName  | string |账户名|

### /api/user/changePassword
修改账号密码

| 参数名       | 参数类型  | 参数说明|
|-----------| ----  |----|
| accountName  | string |账户名|
| password  | string |密码|

### /api/user/login
用户登录

| 参数名       | 参数类型  | 参数说明|
|-----------| ----  |----|
| accountName  | string |账户名|
| password  | string |密码|

### /api/user/logout
用户退出

| 参数名       | 参数类型  | 参数说明|
|-----------| ----  |----|
| accountName  | string |账户名|